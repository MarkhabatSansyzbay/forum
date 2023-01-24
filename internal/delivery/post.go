package delivery

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"forum/internal/models"
	"forum/internal/service"
)

func IDFromURL(url, prefix string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(url, prefix))
}

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(models.User)
	postID, err := IDFromURL(r.URL.Path, "/posts/")
	if err != nil {
		h.errorPage(w, http.StatusNotFound, fmt.Errorf("error getting post ID: %s", err))
		return
	}

	switch r.Method {
	case http.MethodGet:
		post, err := h.services.Post.PostById(postID, user.ID)
		if err != nil {
			if errors.Is(err, service.ErrNoPost) {
				h.errorPage(w, http.StatusNotFound, nil)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		comments, err := h.services.Commentary.CommentsByPostID(postID, user.ID)
		if err != nil {
			log.Printf("error getting comments by post ID: %s", err)
		}

		data := models.TemplateData{
			Template: "post-page",
			User:     user,
			Post:     post,
			Comments: comments,
		}

		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
		}
	case http.MethodPost:
		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		commentContent, ok := r.Form["comment"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		comment := models.Comment{
			UserID:  user.ID,
			PostID:  postID,
			Content: commentContent[0],
		}

		if err := h.services.Commentary.CreateComment(comment); err != nil {
			if errors.Is(err, service.ErrEmptyComment) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data := models.TemplateData{
			Template: "create-post",
			User:     user,
		}

		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		// if err := r.ParseForm(); err != nil {
		// 	h.errorPage(w, http.StatusInternalServerError, err)
		// 	return
		// }

		if err := r.ParseMultipartForm(5 << 20); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
		title, ok1 := r.Form["title"]
		content, ok2 := r.Form["content"]
		category, ok3 := r.Form["category"]

		if !ok1 || !ok2 || !ok3 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		images := r.MultipartForm.File["image"]

		for _, fileHeader := range images {
			if fileHeader.Size > 5<<20 {
				h.errorPage(w, http.StatusBadRequest, errors.New("file size is too big"))
				return
			}
			file, err := fileHeader.Open()
			if err != nil {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "image/svg" {
				http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
				return
			}

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = os.MkdirAll("./uploads", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			defer f.Close()

			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		post := models.Post{
			Title:      title[0],
			AuthorID:   user.ID,
			Content:    content[0],
			Categories: category,
		}

		if err := h.services.Post.CreatePost(post); err != nil {
			if errors.Is(err, service.ErrEmptyPost) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) reactToPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(models.User)

	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	if r.Method == http.MethodGet {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	reaction, ok := r.Form["react"]
	if !ok {
		h.errorPage(w, http.StatusBadRequest, nil)
		return
	}

	id, err := IDFromURL(r.URL.Path, "/posts/react/")
	if err != nil {
		h.errorPage(w, http.StatusNotFound, fmt.Errorf("error getting post ID: %s", err))
		return
	}

	if err := h.services.Reaction.ReactToPost(id, user.ID, reaction[0]); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
}

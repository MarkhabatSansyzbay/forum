package delivery

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	user := r.Context().Value(contextKeyUser).(models.User)

	switch r.Method {
	case http.MethodGet:
		var posts []models.Post
		if len(r.URL.Query()) == 0 {
			var err error
			posts, err = h.services.Post.AllPosts(user.ID)
			if err != nil {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			category := r.URL.Query().Get("category")
			var err error
			posts, err = h.services.Post.PostsByCategory(user.ID, category)
			if err != nil {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
		}

		data := models.TemplateData{
			User:     user,
			Posts:    posts,
			Template: "index",
		}

		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
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

		postID, ok1 := r.Form["postID"]
		react, ok2 := r.Form["react"]

		if !ok1 || !ok2 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		id, err := strconv.Atoi(postID[0])
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		if err := h.services.Reaction.ReactToPost(id, user.ID, react[0]); err != nil {
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

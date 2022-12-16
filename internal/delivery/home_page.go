package delivery

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}

	user := r.Context().Value(contextKeyUser).(models.User)
	// posts, err := h.services.GetAllPosts()
	// if err != nil {
	// 	h.errorPage(w, http.StatusInternalServerError, err)
	// 	return
	// }

	data := models.TemplateData{
		User: user,
		// Posts:    posts,
		Template: "index",
	}

	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}
}

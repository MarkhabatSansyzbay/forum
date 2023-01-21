package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) reactComment(w http.ResponseWriter, r *http.Request) {
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

	react, ok1 := r.Form["react"]
	commentIDVal, ok2 := r.Form["commentID"]

	if !ok1 || !ok2 {
		h.errorPage(w, http.StatusBadRequest, nil)
		return
	}

	user := r.Context().Value(contextKeyUser).(models.User)

	commentID, err := strconv.Atoi(commentIDVal[0])
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	postID, err := h.services.Reaction.ReactToComment(commentID, user.ID, react[0])
	if err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%v", postID), http.StatusSeeOther)
}

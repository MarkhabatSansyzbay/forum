package delivery

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) setSession(w http.ResponseWriter, user *models.User, isOauth2 bool) error {
	session, err := h.services.Authorization.SetSession(user.Username, user.Password, isOauth2)
	if err != nil {
		// if errors.Is(err, service.ErrNoUser) || errors.Is(err, service.ErrWrongPassword) {
		// 	h.errorPage(w, http.StatusUnauthorized, err)
		// 	return
		// }
		// h.errorPage(w, http.StatusInternalServerError, err)
		return err
	}

	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Path:    "/",
		Expires: session.ExpirationDate,
	}
	http.SetCookie(w, cookie)

	return nil
}

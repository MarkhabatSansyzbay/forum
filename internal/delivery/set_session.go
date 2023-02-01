package delivery

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) setSession(w http.ResponseWriter, user *models.User) error {
	session, err := h.services.Authorization.SetSession(user)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpirationDate,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}

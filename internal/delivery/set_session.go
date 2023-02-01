package delivery

import (
	"net/http"

	"forum/internal/models"
)

func (h *Handler) setSession(w http.ResponseWriter, user *models.User, isOauth2 bool) error {
	session, err := h.services.Authorization.SetSession(user.Username, user.Password, isOauth2)
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

package delivery

import (
	"errors"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/service"
)

const (
	templateSignUp = "templates/sign-up.html"
	templateSignIn = "templates/sign-in.html"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := templateExecute(w, templateSignUp, nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		email, ok1 := r.Form["email"]
		username, ok2 := r.Form["username"]
		password, ok3 := r.Form["password"]

		if !ok1 || !ok2 || !ok3 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		user := models.User{
			Email:    email[0],
			Username: username[0],
			Password: password[0],
		}

		_, err := h.services.CreateUser(user)
		if err != nil {
			if err == service.ErrInvalidData {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			} else {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
		}
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := templateExecute(w, templateSignIn, nil); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		username, ok1 := r.Form["username"]
		password, ok2 := r.Form["password"]

		if !ok1 || !ok2 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		possibleUser := models.User{
			Username: username[0],
			Password: password[0],
		}

		user, err := h.services.GetUser(possibleUser)
		if err != nil {
			if errors.Is(err, service.ErrNoUser) {
				h.errorPage(w, http.StatusUnauthorized, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		session, err := h.services.SetSession(user.Id)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Path:    "/",
			Expires: session.ExpirationDate,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusBadRequest, err)
		return
	}
	h.services.DeleteSession(cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

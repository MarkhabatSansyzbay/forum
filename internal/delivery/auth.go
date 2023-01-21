package delivery

import (
	"errors"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/service"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := models.TemplateData{
			Template: "sign-up",
		}
		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
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
		confirm, ok4 := r.Form["confirm_password"]

		if !ok1 || !ok2 || !ok3 || !ok4 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		user := models.User{
			Email:           email[0],
			Username:        username[0],
			Password:        password[0],
			ConfirmPassword: confirm[0],
		}

		if err := h.services.Authorization.CreateUser(user); err != nil {
			if errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrInvalidPassword) ||
				errors.Is(err, service.ErrInvalidUsername) || errors.Is(err, service.ErrUsernameTaken) ||
				errors.Is(err, service.ErrEmailTaken) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			} else {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := models.TemplateData{
			Template: "sign-in",
		}
		if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
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

		session, err := h.services.Authorization.SetSession(username[0], password[0])
		if err != nil {
			if errors.Is(err, service.ErrNoUser) || errors.Is(err, service.ErrWrongPassword) {
				h.errorPage(w, http.StatusUnauthorized, err)
				return
			}
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
	if r.Method == http.MethodGet {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Authorization.DeleteSession(cookie.Value); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

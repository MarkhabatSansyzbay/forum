package delivery

import (
	"context"
	"fmt"
	"net/http"

	"forum/internal/models"
)

const contextKeyUser contextKey = "user"

type contextKey string

func (h *Handler) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		var user models.User
		switch err {
		case http.ErrNoCookie:
			user = models.User{}
		case nil:
			user, err = h.services.Authorization.UserByToken(cookie.Value)
			if err != nil {
				fmt.Printf("user by token: %s\n", err)
			}
		default:
			h.errorPage(w, http.StatusBadRequest, err)
		}
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

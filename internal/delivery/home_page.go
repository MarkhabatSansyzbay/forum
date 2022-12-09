package delivery

import (
	"fmt"
	"net/http"
)

const templateIndex = "templates/index.html"

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodGet {
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}

	user := r.Context().Value(contextKeyUser)
	fmt.Println(user)
	if err := templateExecute(w, templateIndex, user); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
	}
	// надо взять посты и всю информацию о них
	// если пользователь авторизован, надо чтобы у него была возможность создавать посты, комменты и ставить реакцию
	// read a cookie and check if user is authorized
}

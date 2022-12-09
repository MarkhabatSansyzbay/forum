package delivery

import (
	"log"
	"net/http"
)

func (h *Handler) errorPage(w http.ResponseWriter, status int, err error) {
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
}

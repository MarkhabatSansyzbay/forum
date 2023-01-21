package delivery

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) errorPage(w http.ResponseWriter, status int, err error) {
	var msg string = http.StatusText(status)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.Println(err)
		} else if status != http.StatusNotFound {
			msg = err.Error()
		}
	}

	w.WriteHeader(status)

	data := models.TemplateData{
		Error: models.ErrorMsg{
			Status: status,
			Msg:    msg,
		},
	}

	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}
}

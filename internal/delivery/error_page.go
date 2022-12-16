package delivery

import (
	"log"
	"net/http"

	"forum/internal/models"
)

func (h *Handler) errorPage(w http.ResponseWriter, status int, err error) {
	var msg string = http.StatusText(status)
	if err != nil {
		log.Println(err)
		if status != http.StatusInternalServerError {
			msg = err.Error()
		}
	}

	w.WriteHeader(status)

	data := models.TemplateData{
		Template: "error",
		Error: models.ErrorMsg{
			Status: status,
			Msg:    msg,
		},
	}

	if err := h.tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}
}

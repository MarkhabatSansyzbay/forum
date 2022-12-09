package delivery

import (
	"net/http"
	"text/template"

	"forum/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		services: service,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.middleware(h.homePage))
	mux.HandleFunc("/sign-up", h.signUp)
	mux.HandleFunc("/sign-in", h.signIn)
	mux.HandleFunc("/logout", h.logOut)

	return mux
}

func templateExecute(w http.ResponseWriter, path string, data any) error {
	temp, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	if err = temp.Execute(w, data); err != nil {
		return err
	}

	return nil
}

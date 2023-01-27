package delivery

import (
	"html/template"
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	tmpl     *template.Template
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:     template.Must(template.ParseGlob("templates/*.html")),
		services: service,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.middleware(h.homePage))
	mux.HandleFunc("/sign-up", h.signUp)
	mux.HandleFunc("/sign-in", h.signIn)
	mux.HandleFunc("/sign-out", h.logOut)

	mux.HandleFunc("/sign-in/google", h.googleSignIn)
	mux.HandleFunc("/sign-up/google", h.googleSignUp)
	mux.HandleFunc("/sign-in/google/callback", h.signInCallbackFromGoogle)
	mux.HandleFunc("/sign-up/google/callback", h.signUpCallbackFromGoogle)

	mux.HandleFunc("/posts/", h.middleware(h.postPage))
	mux.HandleFunc("/posts/create", h.middleware(h.createPost))
	mux.HandleFunc("/posts/react/", h.middleware(h.reactToPost))
	mux.HandleFunc("/my-posts", h.middleware(h.myPosts))
	mux.HandleFunc("/liked-posts", h.middleware(h.likedPosts))

	mux.HandleFunc("/comment/react/", h.middleware(h.reactComment))

	mux.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates/"))))

	return mux
}

package handlers

import (
	"net/http"

	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/docs"
	"github.com/piquel-fr/api/services/email"
)

type Handler struct {
	AuthService  auth.AuthService
	DocsService  docs.DocsService
	EmailService email.EmailService
}

func (h *Handler) CreateHttpHandler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", h.rootHandler)
	router.Handle("/auth/", http.StripPrefix("/auth", h.CreateAuthHandler()))
	router.Handle("/profile/", http.StripPrefix("/profile", h.CreateProfileHandler()))
	router.Handle("/docs/", http.StripPrefix("/docs", h.CreateDocsHandler()))
	router.Handle("/email/", http.StripPrefix("/email", h.CreateEmailHandler()))

	return router
}

func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Welcome to the Piquel API! Visit the <a href=\"https://piquel.fr/docs\">API</a> for more information."))
	w.WriteHeader(http.StatusOK)
}

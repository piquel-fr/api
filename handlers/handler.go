package handlers

import (
	"net/http"

	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/docs"
)

type Handler struct {
	AuthService auth.AuthService
	DocsService docs.DocsService
}

func (h *Handler) CreateHttpHandler() http.Handler {
	router := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", h.CreateAuthHandler()))
	router.Handle("/profile/", http.StripPrefix("/profile", h.CreateProfileHandler()))
	router.Handle("/docs/", http.StripPrefix("/docs", h.CreateDocsHandler()))

	return router
}

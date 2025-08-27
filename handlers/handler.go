package handlers

import (
	"net/http"

	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/docs"
	gh "github.com/piquel-fr/api/utils/github"
)

type Handler struct {
	AuthService auth.AuthService
	DocsService docs.DocsService

	config   *models.Configuration
	database *repository.Queries
	gh       *gh.GhWrapper
}

func (h *Handler) CreateHttpHandler(config *models.Configuration, database *repository.Queries, gh *gh.GhWrapper) http.Handler {
	h.config = config
	h.database = database
	h.gh = gh

	router := http.NewServeMux()
	router.Handle("/auth/", http.StripPrefix("/auth", h.CreateAuthHandler()))
	router.Handle("/profile/", http.StripPrefix("/profile", h.CreateProfileHandler()))
	router.Handle("/docs/", http.StripPrefix("/docs", h.CreateDocsHandler()))

	return router
}

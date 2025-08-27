package handlers

import (
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

func (h *Handler) CreateAuthHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /{provider}", h.handleProviderLogin)
	handler.HandleFunc("GET /{provider}/{callback}", h.handleAuthCallback)

	handler.Handle("OPTIONS /{provider}", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{provider}/callback", middleware.CreateOptionsHandler("GET"))

	return handler
}

func (h *Handler) handleProviderLogin(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := h.AuthService.GetProvider(providerName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, provider.AuthCodeURL(r.URL.Query().Get("redirectTo")), http.StatusTemporaryRedirect)
}

func (h *Handler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := h.AuthService.GetProvider(providerName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	token, err := provider.GetOAuthConfig().Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	oauthUser, err := provider.FetchUser(r.Context(), token)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	user, err := h.AuthService.GetUser(r.Context(), oauthUser)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	tokenString, err := h.AuthService.GenerateTokenString(user.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	redirectUrl := h.formatRedirectURL(r.URL.Query().Get("state"), tokenString)
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func (h *Handler) formatRedirectURL(redirectTo string, token string) string {
	return fmt.Sprintf("%s?redirectTo=%s&token=%s", h.config.Envs.AuthCallbackUrl, utils.FormatLocalPathString(redirectTo), token)
}

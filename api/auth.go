package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type AuthHandler struct {
	authService auth.AuthService
}

func CreateAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /policy.json", h.policyHandler)
	handler.HandleFunc("GET /{provider}", h.handleProviderLogin)
	handler.HandleFunc("GET /{provider}/callback", h.handleAuthCallback)

	handler.Handle("OPTIONS /{provider}", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{provider}/callback", middleware.CreateOptionsHandler("GET"))

	return handler
}

func (h *AuthHandler) policyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(auth.Policy); err != nil {
		errors.HandleError(w, r, err)
	}
}

func (h *AuthHandler) handleProviderLogin(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := h.authService.GetProvider(providerName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, provider.AuthCodeURL(r.URL.Query().Get("redirectTo")), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := h.authService.GetProvider(providerName)
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

	user, err := h.authService.GetUser(r.Context(), oauthUser)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	tokenString, err := h.authService.GenerateTokenString(user.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	redirectUrl := h.formatRedirectURL(r.URL.Query().Get("state"), tokenString)
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) formatRedirectURL(redirectTo string, token string) string {
	return fmt.Sprintf("%s?redirectTo=%s&token=%s", config.Envs.AuthCallbackUrl, utils.FormatLocalPathString(redirectTo), token)
}

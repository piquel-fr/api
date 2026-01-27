package api

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type AuthHandler struct {
	userService users.UserService
	authService auth.AuthService
}

func CreateAuthHandler(userService users.UserService, authService auth.AuthService) *AuthHandler {
	return &AuthHandler{userService, authService}
}

func (h *AuthHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /{provider}", h.handleProviderLogin)
	handler.HandleFunc("GET /{provider}/callback", h.handleAuthCallback)

	handler.Handle("OPTIONS /{provider}", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{provider}/callback", middleware.CreateOptionsHandler("GET"))

	return handler
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

	user, err := h.userService.GetUserByEmail(r.Context(), oauthUser.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		user, err = h.userService.RegisterUser(r.Context(), oauthUser.Username, oauthUser.Email, oauthUser.Name, oauthUser.Image, auth.RoleDefault)
	}
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	tokenString, err := h.authService.SignToken(h.authService.GenerateToken(user))
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

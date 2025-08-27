package handlers

import (
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/auth/oauth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/middleware"
	"github.com/piquel-fr/api/utils"
)

func CreateAuthHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /{provider}", handleProviderLogin)
	handler.HandleFunc("GET /{provider}/{callback}", handleAuthCallback)

	handler.Handle("OPTIONS /{provider}", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{provider}/callback", middleware.CreateOptionsHandler("GET"))

	return handler
}

func handleProviderLogin(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := oauth.GetProvider(providerName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, provider.AuthCodeURL(r.URL.Query().Get("redirectTo")), http.StatusTemporaryRedirect)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, err := oauth.GetProvider(providerName)
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

	user, err := auth.GetUser(r.Context(), oauthUser)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	tokenString, err := auth.GenerateTokenString(user.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	redirectUrl := formatRedirectURL(r.URL.Query().Get("state"), tokenString)
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func formatRedirectURL(redirectTo string, token string) string {
	return fmt.Sprintf("%s?redirectTo=%s&token=%s", config.Envs.AuthCallbackUrl, utils.FormatLocalPathString(redirectTo), token)
}

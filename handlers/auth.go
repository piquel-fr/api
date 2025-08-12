package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/auth/oauth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/middleware"
)

func CreateAuthHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /logout", handleLogout)
	handler.HandleFunc("GET /{provider}", handleProviderLogin)
	handler.HandleFunc("GET /{provider}/{callback}", handleAuthCallback)

	handler.Handle("OPTIONS /logout", middleware.CreateOptionsHandler("GET"))
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

	user, err := provider.FetchUser(r.Context(), token)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	userId, err := auth.VerifyUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.StoreUserSession(w, r, userId, &oauth.UserSession{Token: token, User: user})
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	redirectUser(w, r, r.URL.Query().Get("state"))
}

// will be removed when moving to Bearer (should be done in frontend)
func handleLogout(w http.ResponseWriter, r *http.Request) {
	err := auth.RemoveUserSession(w, r)
	if err != nil {
		http.Error(w, "Error removing cookies", http.StatusInternalServerError)
		panic(err)
	}

	http.Redirect(w, r, r.URL.Query().Get("redirectTo"), http.StatusTemporaryRedirect)
}

func redirectUser(w http.ResponseWriter, r *http.Request, redirectTo string) {
	redirectTo = fmt.Sprintf("%s/%s", config.Envs.RedirectTo, strings.Trim(redirectTo, "/"))
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
}

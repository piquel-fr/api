package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/middleware"
	"github.com/piquel-fr/api/utils"
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

const RedirectSession = "redirect_to"

func handleProviderLogin(w http.ResponseWriter, r *http.Request) {
	saveRedirectURL(w, r)
	providerName := r.PathValue("provider")
	provider, ok := auth.Providers[providerName]
	if !ok {
		http.Error(w, fmt.Sprintf("provider %s is not valid", providerName), http.StatusBadRequest)
		return
	}

	state := utils.RandString(64)
	session, err := auth.Store.Get(r, "oauth_session")
	if err != nil {
		panic(err)
	}

	session.Values[providerName] = state
	http.Redirect(w, r, provider.Config.AuthCodeURL(state, provider.AuthCodeOptions), http.StatusTemporaryRedirect)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	provider, ok := auth.Providers[providerName]
	if !ok {
		http.Error(w, fmt.Sprintf("provider %s is not valid", providerName), http.StatusBadRequest)
		return
	}

	token, err := provider.Config.Exchange(r.Context(), r.URL.Query().Get("code"))

	user := models.UserSession{
		AccessToken: token.AccessToken,
		ExpiresAt:   token.Expiry,
	}

	err = provider.FetchUser(provider, &user)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	userId, err := auth.VerifyUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.StoreUserSession(w, r, userId, &user)
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	redirectUser(w, r)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	err := auth.RemoveUserSession(w, r)
	if err != nil {
		http.Error(w, "Error removing cookies", http.StatusInternalServerError)
		panic(err)
	}

	redirectURL := getRedirectURL(r)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func getRedirectURL(r *http.Request) string {
	redirectTo := r.URL.Query().Get("redirectTo")
	return fmt.Sprintf("%s/%s", config.Envs.RedirectTo, strings.Trim(redirectTo, "/"))
}

func saveRedirectURL(w http.ResponseWriter, r *http.Request) {
	redirectURL := getRedirectURL(r)

	session, err := auth.Store.Get(r, RedirectSession)
	if err != nil {
		panic(err)
	}

	session.Values["redirectTo"] = redirectURL

	err = session.Save(r, w)
	if err != nil {
		panic(err)
	}
}

func redirectUser(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, RedirectSession)
	if err != nil {
		panic(err)
	}

	redirectURL := session.Values["redirectTo"]
	session.Values["redirectTo"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	if redirectURL == nil || redirectURL == "" {
		redirectURL = config.Envs.RedirectTo
	}

	http.Redirect(w, r, redirectURL.(string), http.StatusTemporaryRedirect)
}

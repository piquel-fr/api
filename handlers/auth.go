package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth/gothic"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/middleware"
	"github.com/piquel-fr/api/services/users"
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

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	_, err = users.VerifyUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

	redirectUser(w, r)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	userId, err := users.VerifyUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.StoreUserSession(w, r, userId, models.UserSessionFromGothUser(&user))
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	redirectUser(w, r)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.RemoveUserSession(w, r)
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

	session, err := gothic.Store.Get(r, RedirectSession)
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
	session, err := gothic.Store.Get(r, RedirectSession)
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

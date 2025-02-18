package handlers

import (
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/markbates/goth/gothic"
)

func HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
    // Save redirect URL to cookies
    // Verify that it is a registered domain (so piquel.fr)

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
    
    // Check if redirect URL is in cookies
    // Otherise just return redirect to main page
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
        http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

    username, err := users.VerifyUser(r.Context(), &user)
    if err != nil {
        http.Error(w, "Error verifying user", http.StatusInternalServerError)
        panic(err)
    }

	err = auth.StoreUserSession(w, r, username, types.UserSessionFromGothUser(&user))
	if err != nil {
        http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

    // Check if redirect URL is in cookies
    // Otherise just return redirect to main page
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
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
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

    // Check if redirect URL is in request query params
    // Otherise just return redirect to main page
}

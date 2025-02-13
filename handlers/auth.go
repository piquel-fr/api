package handlers

import (
	"context"
	"net/http"

	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/users"
	auth_views "github.com/PiquelChips/piquel.fr/views/auth"
	"github.com/markbates/goth/gothic"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
    _, err := auth.GetSessionUser(r)
    if err == errors.ErrorNotAuthenticated {
	    auth_views.Login(users.GetPageData(w, r)).Render(context.Background(), w)
    }
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

    users.VerifyUser(r.Context(), &user)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		panic(err)
	}

    users.VerifyUser(r.Context(), &user)

	err = auth.StoreUserSession(w, r, user)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
    err := gothic.Logout(w, r)
    if err != nil {
        panic(err)
    }

    auth.RemoveUserSession(w, r)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

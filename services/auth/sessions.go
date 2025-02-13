package auth

import (
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func InitCookieStore() {
    store := sessions.NewCookieStore([]byte(config.Envs.CookiesAuthSecret))

    store.MaxAge(178200)
    store.Options.Path = "/"
    store.Options.HttpOnly = false // should be true if http
    store.Options.Secure = true // should be true if https

    log.Printf("[Cookies] Initialized cookie service!\n")

    gothic.Store = store
}

func GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	user := session.Values["user"]
	if user == nil {
		return goth.User{}, errors.ErrorNotAuthenticated
	}
	return user.(goth.User), nil
}

func StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, err := gothic.Store.Get(r, SessionName)
    if err != nil {
        panic(err)
    }

	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	return nil
}

func RemoveUserSession(w http.ResponseWriter, r *http.Request) {
    session , err := gothic.Store.Get(r, SessionName)
    if err != nil {
        panic(err)
    }
    session.Values["user"] = goth.User{}
    session.Options.MaxAge = -1
    session.Save(r, w)
}


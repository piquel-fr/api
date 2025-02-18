package auth

import (
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

const SessionName = "user_session"

func InitCookieStore() {
    store := sessions.NewCookieStore([]byte(config.Envs.CookiesAuthSecret))

    store.MaxAge(178200)
    store.Options.Path = "/"
    store.Options.HttpOnly = false // should be true if http
    store.Options.Secure = true // should be true if https
    store.Options.Domain = config.Envs.OrgDomain

    log.Printf("[Cookies] Initialized cookie service!\n")

    gothic.Store = store
}

func VerifyUserSession(r *http.Request) error {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	user := session.Values["user"]
	if user == nil {
		return errors.ErrorNotAuthenticated
	}
	return nil
}

func StoreUserSession(w http.ResponseWriter, r *http.Request, username string, userSession *types.UserSession) error {
	session, err := gothic.Store.Get(r, SessionName)
    if err != nil {
        return err
    }

	session.Values["username"] = username
    session.Values["session"] = userSession

	err = session.Save(r, w)
	return err
}

func GetUserSession(r *http.Request) (*types.UserSession, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	userSession := session.Values["session"]
	if userSession == nil {
		return nil, errors.ErrorNotAuthenticated
	}
	return userSession.(*types.UserSession), nil
}

func GetUsername(r *http.Request) (string, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return "", err
	}

	username := session.Values["username"]
	if username == "" {
		return "", errors.ErrorNotAuthenticated
	}
	return username.(string), nil
}

func RemoveUserSession(w http.ResponseWriter, r *http.Request) error {
    session , err := gothic.Store.Get(r, SessionName)
    if err != nil {
        return err
    }
    session.Values["username"] = ""
    session.Values["session"] = types.UserSession{}
    session.Options.MaxAge = -1
    session.Save(r, w)
    return nil
}

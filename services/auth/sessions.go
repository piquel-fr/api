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
	store.Options.HttpOnly = false
	store.Options.Secure = true
	store.Options.Domain = config.Envs.Domain

	log.Printf("[Cookies] Initialized cookie service!\n")

	gothic.Store = store
}

func VerifyUserSession(r *http.Request) error {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	user := session.Values["session"]
	if user == nil {
		return errors.ErrorNotAuthenticated
	}
	return nil
}

func StoreUserSession(w http.ResponseWriter, r *http.Request, userId int32, userSession *types.UserSession) error {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	session.Values["userId"] = userId
	session.Values["session"] = userSession

	return session.Save(r, w)
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

func GetUserId(r *http.Request) (int32, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return 0, err
	}

	userId := session.Values["userId"]
	if userId == 0 || userId == nil {
		return 0, errors.ErrorNotAuthenticated
	}
	return userId.(int32), nil
}

func RemoveUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return err
	}
	session.Values["userId"] = 0
	session.Values["session"] = types.UserSession{}
	session.Options.MaxAge = -1
	session.Save(r, w)
	return nil
}

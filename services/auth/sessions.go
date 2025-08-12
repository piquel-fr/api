package auth

import (
	"log"
	"net/http"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/models"
	"github.com/gorilla/sessions"
)

const SessionName = "user_session"

var Store *sessions.CookieStore

func InitCookieStore() {
	Store = sessions.NewCookieStore([]byte(config.Envs.CookiesAuthSecret))

	Store.MaxAge(178200)
	Store.Options.Path = "/"
	Store.Options.HttpOnly = false
	Store.Options.Secure = true
	Store.Options.Domain = config.Envs.Domain

	log.Printf("[Cookies] Initialized cookie service!\n")
}

func VerifyUserSession(r *http.Request) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	user := session.Values["session"]
	if user == nil {
		return errors.ErrorNotAuthenticated
	}
	return nil
}

func StoreUserSession(w http.ResponseWriter, r *http.Request, userId int32, userSession *models.UserSession) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	session.Values["userId"] = userId
	session.Values["session"] = userSession

	return session.Save(r, w)
}

func GetUserSession(r *http.Request) (*models.UserSession, error) {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	userSession := session.Values["session"]
	if userSession == nil {
		return nil, errors.ErrorNotAuthenticated
	}
	return userSession.(*models.UserSession), nil
}

func GetUserId(r *http.Request) (int32, error) {
	session, err := Store.Get(r, SessionName)
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
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}
	session.Values["userId"] = 0
	session.Values["session"] = models.UserSession{}
	session.Options.MaxAge = -1
	session.Save(r, w)
	return nil
}

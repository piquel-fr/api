package models

import (
	"encoding/gob"
	"time"

	"github.com/markbates/goth"
	repository "github.com/piquel-fr/api/database/generated"
)

func Init() {
	gob.Register(UserSession{})
}

type UserProfile struct {
	*repository.User
	Color    string `json:"color"`
	RoleName string `json:"role_name"`
}

func (profile *UserProfile) GetResourceName() string {
	return "user"
}

func (profile *UserProfile) GetOwner() int32 {
	return profile.ID
}

type UserSession struct {
	AccessToken       string    `json:"access_token"`
	AccessTokenSecret string    `json:"access_token_secret"`
	RefreshToken      string    `json:"refresh_token"`
	ExpiresAt         time.Time `json:"expires_at"`
	IDToken           string    `json:"id_token"`
}

func UserSessionFromGothUser(user *goth.User) *UserSession {
	return &UserSession{
		AccessToken:       user.AccessToken,
		AccessTokenSecret: user.AccessTokenSecret,
		RefreshToken:      user.RefreshToken,
		ExpiresAt:         user.ExpiresAt,
		IDToken:           user.IDToken,
	}
}

type DocsInstance repository.DocsInstance

func (docs *DocsInstance) GetResourceName() string {
	return "docs_instance"
}

func (docs *DocsInstance) GetOwner() int32 {
	return docs.OwnerId
}

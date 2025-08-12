package models

import (
	"encoding/gob"
	"time"

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
	AccessToken           string    `json:"access_token"`
	ExpiresAt             time.Time `json:"expires_at"`
	Email, Name, Username string
	Image, Role           string
}

type DocsInstance repository.DocsInstance

func (docs *DocsInstance) GetResourceName() string {
	return "docs_instance"
}

func (docs *DocsInstance) GetOwner() int32 {
	return docs.OwnerId
}

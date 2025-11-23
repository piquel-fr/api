package models

import "github.com/piquel-fr/api/database/repository"

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

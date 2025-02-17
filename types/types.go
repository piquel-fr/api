package types

import (
	repository "github.com/PiquelChips/piquel.fr/database/generated"
)

type UserProfile struct {
	repository.User
	Color     string `json:"color"`
	Group     string `json:"group"`
	GroupName string `json:"group_name"`
}

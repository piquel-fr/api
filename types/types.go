package types

import (
	repository "github.com/PiquelChips/piquel.fr/database/generated"
)

type UserProfile struct {
	User      repository.User `json:"user"`
	UserColor string          `json:"user_color"`
	UserGroup string          `json:"user_group"`
}

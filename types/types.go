package types

import (
	repository "github.com/PiquelChips/piquel.fr/database/generated"
)

type PageData struct {
    Profile UserProfile
}

type UserProfile struct {
    User repository.User
    UserColor string
    UserGroup string
}

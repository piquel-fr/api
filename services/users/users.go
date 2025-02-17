package users

import (
	"context"
	"time"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/utils"
	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
)

func VerifyUser(context context.Context, inUser *goth.User) {
    _, err := database.Queries.GetUserByEmail(context, inUser.Email)
    if err != nil {
        if err == pgx.ErrNoRows {
            registerUser(context, inUser)
            return
        }
        panic(err)
    }
}

func registerUser(context context.Context, inUser *goth.User) {
    params := repository.AddUserParams{}

    params.Email = inUser.Email
    params.Group = "default"
    params.Image = inUser.AvatarURL
    params.Created = time.Now()
    params.Name = inUser.Name

    switch inUser.Provider {
    case "google":
        params.Username = utils.FormatUsername(inUser.Name)
    case "github":
        params.Username = utils.FormatUsername(inUser.NickName)
    }

    err := database.Queries.AddUser(context, params)
    if err != nil {
        panic(err)
    }
}

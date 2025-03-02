package users

import (
	"context"
	"time"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/permissions"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/PiquelChips/piquel.fr/utils"
	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
)

func VerifyUser(context context.Context, inUser *goth.User) (string, error) {
	user, err := database.Queries.GetUserByEmail(context, inUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return registerUser(context, inUser)
		}
		panic(err)
	}

	return user.Username, nil
}

func registerUser(context context.Context, inUser *goth.User) (string, error) {
	params := repository.AddUserParams{}

	params.Email = inUser.Email
	params.Role = "default"
	params.Image = inUser.AvatarURL
	params.CreatedAt = time.Now()
	params.Name = inUser.Name

	switch inUser.Provider {
	case "google":
		params.Username = utils.FormatUsername(inUser.Name)
	case "github":
		params.Username = utils.FormatUsername(inUser.NickName)
	}

	err := database.Queries.AddUser(context, params)
	return params.Username, err
}

func GetProfile(username string) (*types.UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	profile := &types.UserProfile{User: user}

	role := permissions.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

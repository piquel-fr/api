package users

import (
	"context"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/PiquelChips/piquel.fr/utils"
	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
)

func VerifyUser(context context.Context, inUser *goth.User) (int32, error) {
	user, err := database.Queries.GetUserByEmail(context, inUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return registerUser(context, inUser)
		}
		panic(err)
	}

	return user.ID, nil
}

func registerUser(context context.Context, inUser *goth.User) (int32, error) {
	params := repository.AddUserParams{}

	params.Email = inUser.Email
	params.Role = "default"
	params.Image = inUser.AvatarURL
	params.Name = inUser.Name

	switch inUser.Provider {
	case "google":
		params.Username = utils.FormatUsername(inUser.Name)
	case "github":
		params.Username = utils.FormatUsername(inUser.NickName)
	}

	id, err := database.Queries.AddUser(context, params)
	return id, err
}

func GetProfileFromUsername(username string) (*types.UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	profile := &types.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func GetProfileFromUserId(userId int32) (*types.UserProfile, error) {
	user, err := database.Queries.GetUserById(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	profile := &types.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

package users

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/utils"
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

func GetProfileFromUsername(username string) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func GetProfileFromUserId(userId int32) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserById(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func GetUserFromRequest(r *http.Request) (*repository.User, error) {
	userId, err := auth.GetUserId(r)
	if err != nil {
		return nil, err
	}

	user, err := database.Queries.GetUserById(r.Context(), userId)
	return &user, err
}

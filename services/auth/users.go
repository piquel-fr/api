package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth/oauth"
	"github.com/piquel-fr/api/database"
)

func GetUserProfile(context context.Context, inUser *oauth.User) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserByEmail(context, inUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
			user, err = registerUser(context, inUser)
		}
		if err != nil {
			return nil, err
		}
	}

	profile := &models.UserProfile{User: &user}

	role := Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func registerUser(context context.Context, inUser *oauth.User) (repository.User, error) {
	params := repository.AddUserParams{}

	params.Email = inUser.Email
	params.Username = inUser.Username
	params.Role = "default"
	params.Image = inUser.Image
	params.Name = inUser.Name

	user, err := database.Queries.AddUser(context, params)
	return user, err
}

func GetProfileFromUsername(username string) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := Policy.Roles[profile.Role]

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

	role := Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

package profile

import (
	"context"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
)

func GetProfileFromUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func GetProfileFromUserId(ctx context.Context, userId int32) (*models.UserProfile, error) {
	user, err := database.Queries.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := auth.Policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

package auth

import (
	"context"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
)

type UserProfile struct {
	*repository.User
	Color    string `json:"color"`
	RoleName string `json:"role_name"`
}

func (s *realAuthService) GetProfileFromUsername(ctx context.Context, username string) (*UserProfile, error) {
	user, err := database.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{User: &user}

	role := policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func (s *realAuthService) GetProfileFromUserId(ctx context.Context, userId int32) (*UserProfile, error) {
	user, err := database.Queries.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{User: &user}

	role := policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

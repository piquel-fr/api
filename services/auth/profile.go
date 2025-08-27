package auth

import (
	"context"

	"github.com/piquel-fr/api/models"
)

func (s *realAuthService) GetProfileFromUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	user, err := s.database.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := s.policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

func (s *realAuthService) GetProfileFromUserId(ctx context.Context, userId int32) (*models.UserProfile, error) {
	user, err := s.database.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	profile := &models.UserProfile{User: &user}

	role := s.policy.Roles[profile.Role]

	profile.RoleName = role.Name
	profile.Color = role.Color

	return profile, nil
}

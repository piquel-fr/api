package users

import (
	"context"
	"fmt"
	"slices"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/utils"
)

type UserService interface {
	// getting the user
	GetUserById(ctx context.Context, id int32) (*repository.User, error)
	GetUserByUsername(ctx context.Context, username string) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)

	// managing users
	UpdateUser(ctx context.Context, params repository.UpdateUserParams) error
	UpdateUserAdmin(ctx context.Context, params repository.UpdateUserAdminParams) error
	RegisterUser(ctx context.Context, email, username, name, image, role string) error
	DeleteUser(ctx context.Context, id int32) error

	// other
	FormatAndValidateUsername(username string) (string, error)
	ListUsers(ctx context.Context, offset, limit int32) ([]repository.User, error)
}

type realUserService struct{}

func NewRealUserService() *realUserService {
	config.UsernameBlacklist = []string{"self", "users", "admin", "system"} // TODO: add more
	return &realUserService{}
}

func (s *realUserService) GetUserById(ctx context.Context, id int32) (*repository.User, error) {
	return nil, nil
}

func (s *realUserService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	return nil, nil
}

func (s *realUserService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	return nil, nil
}

func (s *realUserService) UpdateUser(ctx context.Context, params repository.UpdateUserParams) error {
	return nil
}

func (s *realUserService) UpdateUserAdmin(ctx context.Context, params repository.UpdateUserAdminParams) error {
	return nil
}

func (s *realUserService) RegisterUser(ctx context.Context, email, username, name, image, role string) error {
	return nil
}

func (s *realUserService) DeleteUser(ctx context.Context, id int32) error {
	return nil
}

func (s *realUserService) ValidateAndFormatUsername(username string) (string, error) {
	username = utils.FormatUsername(username)
	if slices.Contains(config.UsernameBlacklist, username) {
		return "", fmt.Errorf("username %s is not legal", username)
	}
	return username, nil
}

func (s *realUserService) ListUsers(ctx context.Context, offset, limit int32) ([]repository.User, error) {
	if limit > 200 {
		limit = 200
	}
	return database.Queries.ListUsers(ctx, repository.ListUsersParams{Offset: offset, Limit: limit})
}

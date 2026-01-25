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
	GetUsernameBlacklist() []string

	// getting the user
	GetUserById(ctx context.Context, id int32) (*repository.User, error)
	GetUserByUsername(ctx context.Context, username string) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)

	// managing users
	UpdateUser(ctx context.Context, id int32, username, name, image string) error
	UpdateUserAdmin(ctx context.Context, id int32, username, email, name, image, role string) error
	RegisterUser(ctx context.Context, username, email, name, image, role string) (*repository.User, error)
	DeleteUser(ctx context.Context, id int32) error

	// other
	FormatAndValidateUsername(username string) (string, error)
	ListUsers(ctx context.Context, offset, limit int32) ([]repository.User, error)
}

type realUserService struct{}

func NewRealUserService() *realUserService {
	return &realUserService{}
}

func (s *realUserService) GetUserById(ctx context.Context, id int32) (*repository.User, error) {
	user, err := database.Queries.GetUserById(ctx, id)
	return &user, err
}

func (s *realUserService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	user, err := database.Queries.GetUserByUsername(ctx, username)
	return &user, err
}

func (s *realUserService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	user, err := database.Queries.GetUserByEmail(ctx, email)
	return &user, err
}

func (s *realUserService) UpdateUser(ctx context.Context, id int32, username, name, image string) error {
	username, err := s.FormatAndValidateUsername(username)
	if err != nil {
		return err
	}

	params := repository.UpdateUserParams{
		ID:       id,
		Username: username,
		Name:     name,
		Image:    image,
	}
	return database.Queries.UpdateUser(ctx, params)
}

func (s *realUserService) UpdateUserAdmin(ctx context.Context, id int32, username, email, name, image, role string) error {
	username, err := s.FormatAndValidateUsername(username)
	if err != nil {
		return err
	}

	if err := config.Policy.ValidateRole(role); err != nil {
		return err
	}

	params := repository.UpdateUserAdminParams{
		ID:       id,
		Username: username,
		Email:    email,
		Name:     name,
		Image:    image,
		Role:     role,
	}
	return database.Queries.UpdateUserAdmin(ctx, params)
}

func (s *realUserService) RegisterUser(ctx context.Context, username, email, name, image, role string) (*repository.User, error) {
	username, err := s.FormatAndValidateUsername(username)
	if err != nil {
		return nil, err
	}

	if err := config.Policy.ValidateRole(role); err != nil {
		return nil, err
	}

	params := repository.AddUserParams{
		Username: username,
		Email:    email,
		Name:     name,
		Image:    image,
		Role:     role,
	}

	user, err := database.Queries.AddUser(ctx, params)
	return &user, err
}

func (s *realUserService) DeleteUser(ctx context.Context, id int32) error {
	// TODO: delete user
	return nil
}

func (s *realUserService) FormatAndValidateUsername(username string) (string, error) {
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

func (s *realUserService) GetUsernameBlacklist() []string {
	return []string{"self", "root", "users", "admin", "system"} // TODO: add more
}

package users

import (
	"context"

	"github.com/piquel-fr/api/database/repository"
)

type UserService interface {
	// getting the user
	GetUserById(ctx context.Context, id int32) (*repository.User, error)
	GetUserByUsername(ctx context.Context, username string) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)

	// managing users
	UpdateUser(ctx context.Context, params repository.UpdateUserParams) error
	UpdateUserAdmin(ctx context.Context, params repository.UpdateUserAdminParams) error
	RegisterUser(ctx context.Context, email, username, name, image, role string) error // does not validate the role
	DeleteUser(ctx context.Context, id int32) error

	// other
	ValidateUsername(username string) error
	ListUsers(ctx context.Context, offset, limit int) ([]repository.User, error)
}

type realUserService struct{}

func NewRealUserService() *realUserService {
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

// does not validate the role
func (s *realUserService) RegisterUser(ctx context.Context, email, username, name, image, role string) error {
	return nil
}

func (s *realUserService) DeleteUser(ctx context.Context, id int32) error {
	return nil
}

func (s *realUserService) ValidateUsername(username string) error {
	return nil
}

func (s *realUserService) ListUsers(ctx context.Context, offset, limit int) ([]repository.User, error) {
	return nil, nil
}

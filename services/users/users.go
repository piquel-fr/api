package users

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/utils/errors"
)

type UserService interface {
	GetUsernameBlacklist() []string

	// getting the user
	GetUserById(ctx context.Context, id int32) (*repository.User, error)
	GetUserByUsername(ctx context.Context, username string) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	GetUserFromContext(ctx context.Context) (*repository.User, error)

	// managing users
	UpdateUser(ctx context.Context, params repository.UpdateUserParams) error
	UpdateUserAdmin(ctx context.Context, params repository.UpdateUserAdminParams) error
	RegisterUser(ctx context.Context, username, email, name, image, role string) (*repository.User, error)
	DeleteUser(ctx context.Context, user *repository.User) error

	// other
	ListUsers(ctx context.Context, offset, limit int32) ([]*repository.User, error)
}

type realUserService struct{}

func NewRealUserService() *realUserService {
	return &realUserService{}
}

func (s *realUserService) GetUserById(ctx context.Context, id int32) (*repository.User, error) {
	return database.Queries.GetUserById(ctx, id)
}

func (s *realUserService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	return database.Queries.GetUserByUsername(ctx, username)
}

func (s *realUserService) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	return database.Queries.GetUserByEmail(ctx, email)
}

func (s *realUserService) GetUserFromContext(ctx context.Context) (*repository.User, error) {
	user, ok := ctx.Value(config.UserContextKey).(*repository.User)
	if !ok {
		return nil, fmt.Errorf("user is not in context")
	}
	return user, nil
}

func (s *realUserService) UpdateUser(ctx context.Context, params repository.UpdateUserParams) error {
	username, err := s.formatAndValidateUsername(ctx, params.Username, false)
	if err != nil {
		return err
	}
	params.Username = username

	return database.Queries.UpdateUser(ctx, params)
}

func (s *realUserService) UpdateUserAdmin(ctx context.Context, params repository.UpdateUserAdminParams) error {
	username, err := s.formatAndValidateUsername(ctx, params.Username, false)
	if err != nil {
		return err
	}
	params.Username = username

	if err := config.Policy.ValidateRole(params.Role); err != nil {
		return err
	}

	return database.Queries.UpdateUserAdmin(ctx, params)
}

func (s *realUserService) RegisterUser(ctx context.Context, username, email, name, image, role string) (*repository.User, error) {
	username, err := s.formatAndValidateUsername(ctx, username, true)
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

	return database.Queries.AddUser(ctx, params)
}

func (s *realUserService) DeleteUser(ctx context.Context, user *repository.User) error {
	// TODO: delete user
	return nil
}

// @param force: if the validation can fail. When creating a new user through OAuth, user creation cannot fail. We will thus create a random one
func (s *realUserService) formatAndValidateUsername(ctx context.Context, username string, force bool) (string, error) {
	// check if username actually changing
	user, err := s.GetUserFromContext(ctx)
	if err != nil {
		return "", err
	}

	if user.Username == username {
		return username, nil
	}

	// formatting
	random := false
	username = strings.ReplaceAll(strings.ToLower(username), " ", "")

	// blacklist
	if slices.Contains(config.UsernameBlacklist, username) {
		random = true
		if !force {
			return "", errors.NewError(fmt.Sprintf("username %s is not legal", username), http.StatusBadRequest)
		}
	}

	// regex
	matched, err := regexp.MatchString("^[a-z0-9]+$", username)
	if !matched {
		random = true
		if !force {
			return "", errors.NewError(fmt.Sprintf("username %s contains illegal characters. only letters and numbers are allowed", username), http.StatusBadRequest)
		}
	}
	if err != nil {
		random = true
		if !force {
			return "", fmt.Errorf("error matching regex in username validation: %w", err)
		}
	}

	// already existing users
	names, err := database.Queries.ListUserNames(ctx)
	if err != nil {
		random = true
		if !force {
			return "", err
		}
	}

	if slices.Contains(names, username) {
		random = true
		if !force {
			return "", errors.NewError(fmt.Sprintf("username %s is already taken", username), http.StatusBadRequest)
		}
	}

	// random generation
	if random {
		username = rand.Text()
		username, err = s.formatAndValidateUsername(ctx, username, true)
		if err != nil {
			return "", fmt.Errorf("something terrible happened in username validation:\n\tusername: %s\n\terror: %w", username, err)
		}
	}

	return username, nil
}

func (s *realUserService) ListUsers(ctx context.Context, offset, limit int32) ([]*repository.User, error) {
	if limit > config.MaxLimit {
		limit = config.MaxLimit
	}
	return database.Queries.ListUsers(ctx, offset, limit)
}

func (s *realUserService) GetUsernameBlacklist() []string {
	return []string{"self", "root", "users", "admin", "system"} // TODO: add more
}

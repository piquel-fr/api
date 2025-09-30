package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth/oauth"
	"github.com/piquel-fr/api/utils/errors"
)

type AuthService interface {
	GenerateTokenString(userId int32) (string, error)
	GetToken(r *http.Request) (*jwt.Token, error)
	GetUserId(r *http.Request) (int32, error)
	GetUserFromRequest(r *http.Request) (*repository.User, error)
	GetUser(ctx context.Context, inUser *oauth.User) (*repository.User, error)
	GetUserFromUsername(ctx context.Context, username string) (repository.User, error)
	Authorize(request *Request) error
	GetProvider(name string) (oauth.Provider, error)

	GetProfileFromUsername(ctx context.Context, username string) (*models.UserProfile, error)
	GetProfileFromUserId(ctx context.Context, userId int32) (*models.UserProfile, error)
}

type realAuthService struct {
	providers map[string]oauth.Provider
}

func NewRealAuthService() *realAuthService {
	service := &realAuthService{providers: oauth.GetProviders()}
	return service
}

func (s *realAuthService) GenerateTokenString(userId int32) (string, error) {
	idString := strconv.Itoa(int(userId))
	token := jwt.NewWithClaims(config.JWTSigningMethod,
		jwt.RegisteredClaims{
			Subject: idString,
		})

	return token.SignedString(config.Envs.JWTSigningSecret)
}

func (s *realAuthService) GetToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ""

	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.ErrorNotAuthenticated
	}
	tokenString = parts[1]

	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return config.Envs.JWTSigningSecret, nil
	})
}

func (s *realAuthService) GetUserId(r *http.Request) (int32, error) {
	token, err := s.GetToken(r)
	if err != nil {
		return 0, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}

func (s *realAuthService) GetUserFromRequest(r *http.Request) (*repository.User, error) {
	userId, err := s.GetUserId(r)
	if err != nil {
		return nil, err
	}

	user, err := database.Queries.GetUserById(r.Context(), userId)
	return &user, err
}

func (s *realAuthService) GetUser(ctx context.Context, inUser *oauth.User) (*repository.User, error) {
	user, err := database.Queries.GetUserByEmail(ctx, inUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return s.registerUser(ctx, inUser)
		}
		return nil, err
	}

	return &user, nil
}

func (s *realAuthService) registerUser(ctx context.Context, inUser *oauth.User) (*repository.User, error) {
	params := repository.AddUserParams{}

	params.Email = inUser.Email
	params.Username = inUser.Username
	params.Role = "default"
	params.Image = inUser.Image
	params.Name = inUser.Name

	user, err := database.Queries.AddUser(ctx, params)
	return &user, err
}

func (s *realAuthService) GetUserFromUsername(ctx context.Context, username string) (repository.User, error) {
	return database.Queries.GetUserByUsername(ctx, username)
}

func (s *realAuthService) GetProvider(name string) (oauth.Provider, error) {
	provider, ok := s.providers[name]
	if !ok {
		return nil, errors.NewError(fmt.Sprintf("provider %s does not exist", name), http.StatusBadRequest)
	}
	return provider, nil
}

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/oauth"
)

type AuthService interface {
	GetPolicy() *config.PolicyConfiguration
	GetProvider(name string) (oauth.Provider, error)

	// token management
	GenerateToken(user *repository.User) *jwt.Token // TODO: also save expiry and refresh
	SignToken(token *jwt.Token) (string, error)
	getTokenFromRequest(r *http.Request) (*jwt.Token, error)
	getUserFromToken(ctx context.Context, token *jwt.Token) (*repository.User, error)

	// authorization
	Authorize(request *config.AuthRequest) error
	AuthMiddleware(next http.Handler) http.Handler
}

type realAuthService struct {
	userService users.UserService
}

func NewRealAuthService(userService users.UserService) *realAuthService {
	return &realAuthService{userService}
}

func (s *realAuthService) GetPolicy() *config.PolicyConfiguration { return &policy }

func (s *realAuthService) GetProvider(name string) (oauth.Provider, error) {
	provider, ok := oauth.Providers[name]
	if !ok {
		return nil, errors.NewError(fmt.Sprintf("provider %s does not exist", name), http.StatusBadRequest)
	}
	return provider, nil
}

// TODO: also save expiry and refresh
func (s *realAuthService) GenerateToken(user *repository.User) *jwt.Token {
	idString := strconv.Itoa(int(user.ID))
	return jwt.NewWithClaims(config.JWTSigningMethod,
		jwt.RegisteredClaims{
			Subject: idString,
		})
}

func (s *realAuthService) SignToken(token *jwt.Token) (string, error) {
	return token.SignedString(config.Envs.JWTSigningSecret)
}

func (s *realAuthService) getTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.ErrorNotAuthenticated
	}
	tokenString := parts[1]

	return jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return config.Envs.JWTSigningSecret, nil
	})
}

func (s *realAuthService) getUserFromToken(ctx context.Context, token *jwt.Token) (*repository.User, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return nil, err
	}

	return s.userService.GetUserById(ctx, int32(id))
}

func (s *realAuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		token, err := s.getTokenFromRequest(r)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		user, err := s.getUserFromToken(r.Context(), token)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		newReq := r.WithContext(context.WithValue(r.Context(), config.UserContextKey, user))
		next.ServeHTTP(w, newReq)
	})
}

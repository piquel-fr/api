package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/oauth"
)

type JwtClaims struct {
	User *repository.User `json:"user"`
	jwt.RegisteredClaims
}

type AuthService interface {
	GetPolicy() *config.PolicyConfiguration
	GetProvider(name string) (oauth.Provider, error)

	// token management
	FinishAuth(user *repository.User, w http.ResponseWriter) error // sets the users refresh & access tokens
	Refresh(w http.ResponseWriter, r *http.Request) error          // refreshes the user's tokens

	// authorization
	Authorize(request *config.AuthRequest) error
	AuthMiddleware(next http.Handler) http.Handler
}

type realAuthService struct {
	userService users.UserService
}

func NewRealAuthService(userService users.UserService) AuthService {
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

func (s *realAuthService) FinishAuth(user *repository.User, w http.ResponseWriter) error {
	// TODO
	// 1. create session
	// 2. save session to DB
	// 3. generate refresh_token
	// 4. generate access_token
	// 5. write access_token & refresh_token to cookies
	return nil
}

func (s *realAuthService) Refresh(w http.ResponseWriter, r *http.Request) error {
	// TODO
	// 1. calculate the refresh token hash
	// 2. verify with the stored information in DB
	// 3. generate new refresh_token
	// 4. generate new access_token
	// 5. update the DB session
	// 6. write access_token & refresh_token to cookies
	return nil
}

// TODO: also save expiry and refresh
func (s *realAuthService) generateAccessToken(user *repository.User) *jwt.Token {
	idString := strconv.Itoa(int(user.ID))
	return jwt.NewWithClaims(config.JWTSigningMethod,
		JwtClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: idString,
			},
		})
}

func (s *realAuthService) generateRefreshTokenHash(token string, userId int32, expiry time.Time) (string, error) {
	// TODO: hash all of this
	return "", nil
}

func (s *realAuthService) signToken(token *jwt.Token) (string, error) {
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

package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils"
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
	// 1. generate refresh_token
	// 2. create session & save to DB
	// 3. generate access_token
	// 4. write access_token & refresh_token to cookies
	return nil
}

func (s *realAuthService) Refresh(w http.ResponseWriter, r *http.Request) error {
	// TODO
	// 1. hash refresh_token
	// 2. get session from hash (return 404 if not in DB)
	// 3. verify expiry
	// 4. generate new refresh_token
	// 5. generate new access_token
	// 6. update the DB session (update HASH & push back expiry
	// 7. write access_token & refresh_token to cookies
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

// returns token, hash
func (s *realAuthService) generateRefreshToken() (string, string) {
	token := utils.GenerateSecureToken(32)
	return token, s.hashRefreshToken(token)
}

func (s *realAuthService) verifyRefreshToken(token string, hash string) bool {
	return s.hashRefreshToken(token) == hash
}

func (s *realAuthService) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
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

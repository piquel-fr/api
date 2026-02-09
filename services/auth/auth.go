package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/oauth"
)

const (
	refreshKey = "refresh_token"
	accessKey  = "access_token"
)

type JwtClaims struct {
	User *repository.User `json:"user"`
	jwt.RegisteredClaims
}

type AuthService interface {
	GetPolicy() *config.PolicyConfiguration
	GetProvider(name string) (oauth.Provider, error)

	// token management
	FinishAuth(user *repository.User, r *http.Request, w http.ResponseWriter) error // sets the users refresh & access tokens
	Refresh(w http.ResponseWriter, r *http.Request) error                           // refreshes the user's tokens
	Logout(w http.ResponseWriter, r *http.Request) error

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

func (s *realAuthService) FinishAuth(user *repository.User, r *http.Request, w http.ResponseWriter) error {
	ipAddress := strings.Split(r.RemoteAddr, ":")[0]

	refreshExpiry := time.Hour * 24 * 30 // 30 days
	refreshToken, refreshHash := s.generateRefreshToken(ipAddress)

	accessExpiry := time.Minute * 5 // 5 minutes
	accessToken := s.generateAccessToken(user, time.Now().Add(accessExpiry))
	accessTokenString, err := s.signToken(accessToken)
	if err != nil {
		return err
	}

	sessionParams := repository.AddSessionParams{
		UserId:    user.ID,
		TokenHash: refreshHash,
		UserAgent: r.Header.Get("User-Agent"),
		IpAdress:  ipAddress,
		ExpiresAt: time.Now().Add(refreshExpiry), // one month
	}
	if _, err := database.Queries.AddSession(r.Context(), sessionParams); err != nil {
		return err
	}

	w.Header().Add("Set-Cookie", utils.GenerateSetCookie(refreshKey, refreshToken, config.Envs.Domain, "/auth/refresh", "Strict", refreshExpiry))
	w.Header().Add("Set-Cookie", utils.GenerateSetCookie(accessKey, accessTokenString, config.Envs.Domain, "/", "Lax", accessExpiry))
	return nil
}

func (s *realAuthService) Refresh(w http.ResponseWriter, r *http.Request) error {
	ipAddress := strings.Split(r.RemoteAddr, ":")[0]
	cookies := utils.GetCookiesFromStr(r.Header.Get("Cookie"))

	hash := s.hashRefreshToken(cookies[refreshKey], ipAddress)
	session, err := database.Queries.GetSessionFromHash(r.Context(), hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.ErrorNotAuthenticated
	}
	if err != nil {
		return err
	}

	log.Printf("here")
	if time.Now().After(session.ExpiresAt) {
		return errors.ErrorNotAuthenticated
	}
	log.Printf("expired")
	if !s.verifyRefreshToken(cookies[refreshKey], ipAddress, hash) {
		return errors.ErrorNotAuthenticated
	}
	log.Printf("hash")

	refreshExpiry := time.Hour * 24 * 30 // 30 days
	refreshToken, refreshHash := s.generateRefreshToken(ipAddress)

	user, err := s.userService.GetUserById(r.Context(), session.UserId)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.ErrorNotAuthenticated
	}
	if err != nil {
		return err
	}
	accessExpiry := time.Minute * 5 // 5 minutes
	accessToken := s.generateAccessToken(user, time.Now().Add(accessExpiry))
	accessTokenString, err := s.signToken(accessToken)
	if err != nil {
		return err
	}

	updateSessionParams := repository.UpdateSessionParams{
		UserId:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(refreshExpiry),
	}
	if err := database.Queries.UpdateSession(r.Context(), updateSessionParams); err != nil {
		return err
	}

	w.Header().Add("Set-Cookie", utils.GenerateSetCookie(refreshKey, refreshToken, config.Envs.Domain, "/auth/refresh", "Strict", refreshExpiry))
	w.Header().Add("Set-Cookie", utils.GenerateSetCookie(accessKey, accessTokenString, config.Envs.Domain, "/", "Lax", accessExpiry))
	return nil
}

func (s *realAuthService) Logout(w http.ResponseWriter, r *http.Request) error {
	// TODO
	// 1. hash the token
	// 2. clear the token from the database

	w.Header().Add("Set-Cookie", utils.GenerateClearCookie(refreshKey, config.Envs.Domain, "/auth/refresh"))
	w.Header().Add("Set-Cookie", utils.GenerateClearCookie(accessKey, config.Envs.Domain, "/"))
	return nil
}

func (s *realAuthService) generateAccessToken(user *repository.User, expiresAt time.Time) *jwt.Token {
	idString := strconv.Itoa(int(user.ID))
	return jwt.NewWithClaims(config.JWTSigningMethod,
		JwtClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   idString,
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
}

// returns token, hash
func (s *realAuthService) generateRefreshToken(ipAdress string) (string, string) {
	token := utils.GenerateSecureToken(32)
	return token, s.hashRefreshToken(token, ipAdress)
}

func (s *realAuthService) verifyRefreshToken(token, ipAdress, hash string) bool {
	return s.hashRefreshToken(token, ipAdress) == hash
}

func (s *realAuthService) hashRefreshToken(token, ipAdress string) string {
	hash := sha256.Sum256(fmt.Appendf([]byte(token), "-%s", ipAdress))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (s *realAuthService) signToken(token *jwt.Token) (string, error) {
	return token.SignedString(config.Envs.JWTSigningSecret)
}

func (s *realAuthService) getTokenFromRequest(r *http.Request) (*jwt.Token, *JwtClaims, error) {
	cookies := utils.GetCookiesFromStr(r.Header.Get("Cookie"))
	tokenString := cookies[accessKey]

	claims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return config.Envs.JWTSigningSecret, nil
	})

	return token, claims, err
}

func (s *realAuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		token, claims, err := s.getTokenFromRequest(r)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		if expiration, err := token.Claims.GetExpirationTime(); err != nil {
			errors.HandleError(w, r, err)
			return
		} else if time.Now().After(expiration.Time) {
			errors.HandleError(w, r, errors.ErrorNotAuthenticated)
			return
		}

		newReq := r.WithContext(context.WithValue(r.Context(), config.UserContextKey, claims.User))
		next.ServeHTTP(w, newReq)
	})
}

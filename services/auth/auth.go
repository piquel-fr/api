package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth/oauth"
	"github.com/piquel-fr/api/services/config"
)

func InitAuthService() {
	oauth.InitOAuth()
}

func GenerateTokenString(userId int32) (string, error) {
	idString := strconv.Itoa(int(userId))
	token := jwt.NewWithClaims(config.Configuration.JWTSigningMethod,
		jwt.RegisteredClaims{
			Subject: idString,
		})

	return token.SignedString(config.Envs.JWTSigningSecret)
}

func GetToken(r *http.Request) (*jwt.Token, error) {
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

func GetUserId(r *http.Request) (int32, error) {
	token, err := GetToken(r)
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

func GetUserFromRequest(r *http.Request) (*repository.User, error) {
	userId, err := GetUserId(r)
	if err != nil {
		return nil, err
	}

	user, err := database.Queries.GetUserById(r.Context(), userId)
	return &user, err
}

func GetUser(ctx context.Context, inUser *oauth.User) (*repository.User, error) {
	user, err := database.Queries.GetUserByEmail(ctx, inUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return registerUser(ctx, inUser)
		}
		return nil, err
	}

	return &user, nil
}

func registerUser(ctx context.Context, inUser *oauth.User) (*repository.User, error) {
	params := repository.AddUserParams{}

	params.Email = inUser.Email
	params.Username = inUser.Username
	params.Role = "default"
	params.Image = inUser.Image
	params.Name = inUser.Name

	user, err := database.Queries.AddUser(ctx, params)
	return &user, err
}

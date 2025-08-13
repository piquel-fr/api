package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	if len(parts) == 2 && parts[0] == "Bearer" {
		tokenString = parts[1]
	}

	if tokenString == "" {
		cookie := r.Header.Get("Cookie")
		tokenString = strings.TrimPrefix(cookie, "jwt=")
	}

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

package config

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/piquel-fr/api/models"
)

var Envs models.EnvsConfig
var Configuration = models.Configuration{
	MaxDocsInstanceCount: 3,
	JWTSigningMethod:     jwt.SigningMethodHS256,
}

func LoadConfig() {
	godotenv.Load()
	log.Printf("[Config] Loading configuration...")

	// Load config from environment
	Envs = models.EnvsConfig{
		AuthCallbackUrl:    getEnv("AUTH_CALLBACK"),
		Url:                getEnv("URL"),
		Port:               getDefaultEnv("PORT", "80"),
		DBURL:              getEnv("DB_URL"),
		GoogleClientID:     getEnv("AUTH_GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("AUTH_GOOGLE_CLIENT_SECRET"),
		GithubClientID:     getEnv("AUTH_GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnv("AUTH_GITHUB_CLIENT_SECRET"),
		GithubApiToken:     getEnv("GITHUB_API_TOKEN"),
		JWTSigningSecret:   []byte(getEnv("JWT_SECRET")),
	}

	log.Printf("[Config] Loaded environment configuration!")
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Fatalf("Environment variable %s is not set", key)
	return ""
}

func getDefaultEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

package config

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type EnvsConfig struct {
	AuthCallbackUrl string
	Url             string
	Port            string
	DBURL           string
	GithubApiToken  string

	// auth
	JWTSigningSecret   []byte
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string

	// mail
	SmtpHost string
	ImapHost string
	ImapPort string
}

var Envs EnvsConfig
var MaxDocsInstanceCount int64 = 3
var JWTSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256

func LoadConfig() {
	godotenv.Load()
	log.Printf("[Config] Loading configuration...")

	// Load config from environment
	Envs = EnvsConfig{
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
		SmtpHost:           getEnv("SMTP_HOST"),
		ImapHost:           getEnv("IMAP_HOST"),
		ImapPort:           getDefaultEnv("IMAP_PORT", "993"),
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

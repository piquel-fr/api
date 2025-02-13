package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Envs Config

func LoadConfig() {
	godotenv.Load()
	log.Printf("[Config] Loading configuration...")

    // Load config from environment
	Envs = Config{
		PublicHost:         getEnv("PUBLIC_HOST"),
		Host:               getEnv("HOST"),
		Port:               getEnv("PORT"),
		SSL:                getEnv("SSL"),
		DB_URL:             getEnv("DB_URL"),
		CookiesAuthSecret:  getEnv("COOKIES_AUTH_SECRET"),
		GoogleClientID:     getEnv("AUTH_GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("AUTH_GOOGLE_CLIENT_SECRET"),
		GithubClientID:     getEnv("AUTH_GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnv("AUTH_GITHUB_CLIENT_SECRET"),
		ConfigPath:         getDefaultEnv("CONFIG_PATH", "config"),
	}
	log.Printf("[Config] Loaded environment configuration!")

    loadRouteConfig()
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

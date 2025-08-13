package models

import "github.com/golang-jwt/jwt/v5"

type EnvsConfig struct {
	AuthCallbackUrl    string
	Url                string
	Port               string
	DBURL              string
	GithubApiToken     string
	JWTSigningSecret   []byte
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
}

type Configuration struct {
	MaxDocsInstanceCount int64
	JWTSigningMethod     jwt.SigningMethod
}

package auth

import (
	"fmt"
	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/piquel-fr/api/services/config"
)

func InitAuthentication() {
	goth.UseProviders(
		google.New(
			config.Envs.GoogleClientID,
			config.Envs.GoogleClientSecret,
			buildCallbackURL("google"),
			"email", "profile",
		),
		github.New(
			config.Envs.GithubClientID,
			config.Envs.GithubClientSecret,
			buildCallbackURL("github"),
			"user:email",
		),
	)

	log.Printf("[Auth] Initialized auth service!\n")
}

func buildCallbackURL(provider string) string {
	var url string
	if config.Envs.SSL == "true" {
		url = fmt.Sprintf("https://%s/auth/%s/callback", config.Envs.Host, provider)
	} else {
		url = fmt.Sprintf("http://%s:%s/auth/%s/callback", config.Envs.Host, config.Envs.Port, provider)
	}
	log.Printf("[Auth] Added auth provider listener for %s on %s", provider, url)
	return url
}

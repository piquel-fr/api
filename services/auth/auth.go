package auth

import (
	"fmt"
	"log"

	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const SessionName = "user_session"

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
		url = fmt.Sprintf("https://%s/auth/%s/callback", config.Envs.PublicHost, provider)
	} else {
		url = fmt.Sprintf("http://%s:%s/auth/%s/callback", config.Envs.PublicHost, config.Envs.Port, provider)
	}
	log.Printf("[Auth] Added auth provider listener for %s on %s", provider, url)
	return url
}

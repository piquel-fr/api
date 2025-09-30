package oauth

import (
	"context"
	"fmt"

	"github.com/piquel-fr/api/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Provider interface {
	GetOAuthConfig() *oauth2.Config
	FetchUser(context.Context, *oauth2.Token) (*User, error)
	AuthCodeURL(state string) string
}

type User struct{ Name, Username, Email, Image string }

func GetProviders() map[string]Provider {
	return map[string]Provider{
		"github": &github{
			config: oauth2.Config{
				ClientID:     config.Envs.GithubClientID,
				ClientSecret: config.Envs.GithubClientSecret,
				RedirectURL:  buildCallbackURL(config.Envs.Url, "github"),
				Scopes:       []string{"user:email"},
				Endpoint:     endpoints.GitHub,
			},
		},
		"google": &google{
			config: oauth2.Config{
				ClientID:     config.Envs.GoogleClientID,
				ClientSecret: config.Envs.GoogleClientSecret,
				RedirectURL:  buildCallbackURL(config.Envs.Url, "google"),
				Scopes:       []string{"email", "profile"},
				Endpoint:     endpoints.Google,
			},
			authCodeOptions: oauth2.AccessTypeOffline,
		},
	}
}

func buildCallbackURL(url, provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", url, provider)
}

package oauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/utils/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Provider interface {
	GetOAuthConfig() *oauth2.Config
	FetchUser(context.Context, *oauth2.Token) (*User, error)
	AuthCodeURL(state string) string
}

type User struct{ Name, Username, Email, Image string }

var providers map[string]Provider

func InitOAuth() {
	providers = map[string]Provider{
		"github": &github{
			config: oauth2.Config{
				ClientID:     config.Envs.GithubClientID,
				ClientSecret: config.Envs.GithubClientSecret,
				RedirectURL:  buildCallbackURL("github"),
				Scopes:       []string{"user:email"},
				Endpoint:     endpoints.GitHub,
			},
		},
		"google": &google{
			config: oauth2.Config{
				ClientID:     config.Envs.GoogleClientID,
				ClientSecret: config.Envs.GoogleClientSecret,
				RedirectURL:  buildCallbackURL("google"),
				Scopes:       []string{"email", "profile"},
				Endpoint:     endpoints.Google,
			},
			authCodeOptions: oauth2.AccessTypeOffline,
		},
	}
}

func GetProvider(name string) (Provider, error) {
	provider, ok := providers[name]
	if !ok {
		return nil, errors.NewError(fmt.Sprintf("provider %s does not exist", name), http.StatusBadRequest)
	}
	return provider, nil
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", config.Envs.Url, provider)
}

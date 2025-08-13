package oauth

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Provider interface {
	GetName() string
	GetOAuthConfig() *oauth2.Config
	FetchUser(context.Context, *oauth2.Token) (*User, error)
	AuthCodeURL(state string) string
}

type User struct{ Name, Username, Email, Image string }

type UserSession struct {
	*oauth2.Token
	*User
}

var providers map[string]Provider

func InitOAuth() {
	gob.Register(UserSession{})
	providers = map[string]Provider{
		"github": &github{
			name: "github",
			config: oauth2.Config{
				ClientID:     config.Envs.GithubClientID,
				ClientSecret: config.Envs.GithubClientSecret,
				RedirectURL:  buildCallbackURL("github"),
				Scopes:       []string{"user:email"},
				Endpoint:     endpoints.GitHub,
			},
		},
		"google": &google{
			name: "google",
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
	var url string
	if config.Envs.SSL == "true" {
		url = fmt.Sprintf("https://%s/auth/%s/callback", config.Envs.Host, provider)
	} else {
		url = fmt.Sprintf("http://%s:%s/auth/%s/callback", config.Envs.Host, config.Envs.Port, provider)
	}
	return url
}

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type Provider struct {
	Name            string
	Config          oauth2.Config
	AuthCodeOptions oauth2.AuthCodeOption
	FetchUser       func(Provider, *models.UserSession) error
}

var Providers map[string]Provider

func InitAuthService() {
	Providers = map[string]Provider{
		"github": {
			Name: "github",
			Config: oauth2.Config{
				ClientID:     config.Envs.GithubClientID,
				ClientSecret: config.Envs.GithubClientSecret,
				RedirectURL:  buildCallbackURL("github"),
				Scopes:       []string{"user:email"},
				Endpoint:     endpoints.GitHub,
			},
			FetchUser: func(provider Provider, user *models.UserSession) error {
				jwt.New(jwt.SigningMethodES256)
				const profileURL = "https://api.github.com/user"

				req, err := http.NewRequest("GET", profileURL, nil)
				if err != nil {
					return err
				}

				response, err := provider.Config.Client(context.Background(), &oauth2.Token{AccessToken: user.AccessToken}).Do(req)
				if err != nil {
					return err
				}
				defer response.Body.Close()

				if response.StatusCode != http.StatusOK {
					return fmt.Errorf("GitHub API responded with a %d trying to fetch user information", response.StatusCode)
				}

				u := struct {
					Email   string `json:"email"`
					Name    string `json:"name"`
					Login   string `json:"login"`
					Picture string `json:"avatar_url"`
				}{}

				err = json.NewDecoder(response.Body).Decode(&u)
				if err != nil {
					return err
				}

				user.Name = u.Name
				user.Email = u.Email
				user.Image = u.Picture
				user.Username = utils.FormatUsername(u.Login)

				if user.Email == "" {
					for _, scope := range provider.Config.Scopes {
						if strings.TrimSpace(scope) == "user" || strings.TrimSpace(scope) == "user:email" {
							user.Email, err = getPrivateGithubMail(provider, user)
							if err != nil {
								return err
							}
							break
						}
					}
				}
				return err
			},
			AuthCodeOptions: oauth2.AccessTypeOnline,
		},
		"google": {
			Name: "google",
			Config: oauth2.Config{
				ClientID:     config.Envs.GoogleClientID,
				ClientSecret: config.Envs.GoogleClientSecret,
				RedirectURL:  buildCallbackURL("google"),
				Scopes:       []string{"email", "profile"},
				Endpoint:     endpoints.Google,
			},
			FetchUser: func(provider Provider, user *models.UserSession) error {
				const endpointProfile = "https://www.googleapis.com/oauth2/v2/userinfo"

				response, err := provider.Config.Client(context.Background(), &oauth2.Token{AccessToken: user.AccessToken}).Get(endpointProfile + "?access_token=" + url.QueryEscape(user.AccessToken))
				if err != nil {
					return err
				}
				defer response.Body.Close()

				if response.StatusCode != http.StatusOK {
					return fmt.Errorf("google responded with a %d trying to fetch user information", response.StatusCode)
				}

				responseBytes, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				u := struct {
					Email   string `json:"email"`
					Name    string `json:"name"`
					Picture string `json:"picture"`
				}{}

				if err := json.Unmarshal(responseBytes, &u); err != nil {
					return err
				}

				// Extract the user data we got from Google into our goth.User.
				user.Name = u.Name
				user.Email = u.Email
				user.Image = u.Picture
				user.Username = utils.FormatUsername(u.Name)

				return nil
			},
			AuthCodeOptions: oauth2.AccessTypeOffline,
		},
	}
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

func getPrivateGithubMail(provider Provider, user *models.UserSession) (email string, err error) {
	const emailURL = "https://api.github.com/user/emails"

	req, err := http.NewRequest("GET", emailURL, nil)
	req.Header.Add("Authorization", "Bearer "+user.AccessToken)
	response, err := provider.Config.Client(context.Background(), &oauth2.Token{AccessToken: user.AccessToken}).Do(req)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return email, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return email, fmt.Errorf("GitHub API responded with a %d trying to fetch user email", response.StatusCode)
	}

	var mailList []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	err = json.NewDecoder(response.Body).Decode(&mailList)
	if err != nil {
		return email, err
	}
	for _, v := range mailList {
		if v.Primary && v.Verified {
			return v.Email, nil
		}
	}
	return email, errors.NewError("unable to find a primary github email", http.StatusBadRequest)
}

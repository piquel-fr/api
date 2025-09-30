package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	"golang.org/x/oauth2"
)

const (
	profileURL = "https://api.github.com/user"
	emailURL   = "https://api.github.com/user/emails"
)

type github struct {
	config oauth2.Config
}

func (gh *github) GetOAuthConfig() *oauth2.Config  { return &gh.config }
func (gh *github) AuthCodeURL(state string) string { return gh.config.AuthCodeURL(state) }

func (gh *github) FetchUser(context context.Context, token *oauth2.Token) (*User, error) {
	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := gh.config.Client(context, token).Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API responded with a %d trying to fetch user information", response.StatusCode)
	}

	u := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Login   string `json:"login"`
		Picture string `json:"avatar_url"`
	}{}

	if err = json.NewDecoder(response.Body).Decode(&u); err != nil {
		return nil, err
	}

	user := &User{
		Name:     u.Name,
		Email:    u.Email,
		Image:    u.Picture,
		Username: utils.FormatUsername(u.Login),
	}

	if user.Email == "" {
		for _, scope := range gh.config.Scopes {
			if strings.TrimSpace(scope) == "user" || strings.TrimSpace(scope) == "user:email" {
				user.Email, err = gh.getPrivateEmail(context, token)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	return user, nil
}

func (gh *github) getPrivateEmail(context context.Context, token *oauth2.Token) (string, error) {
	req, err := http.NewRequest("GET", emailURL, nil)
	response, err := gh.config.Client(context, token).Do(req)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.NewError(fmt.Sprintf("GitHub API responded with a %d trying to fetch user email", response.StatusCode), response.StatusCode)
	}

	var mailList []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err = json.NewDecoder(response.Body).Decode(&mailList); err != nil {
		return "", err
	}
	for _, v := range mailList {
		if v.Primary && v.Verified {
			return v.Email, nil
		}
	}
	return "", errors.NewError("unable to find a primary github email", http.StatusBadRequest)
}

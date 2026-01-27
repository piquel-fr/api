package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type google struct {
	config          oauth2.Config
	authCodeOptions oauth2.AuthCodeOption
}

func (gg *google) GetOAuthConfig() *oauth2.Config { return &gg.config }
func (gg *google) AuthCodeURL(state string) string {
	return gg.config.AuthCodeURL(state, gg.authCodeOptions)
}

func (gg *google) FetchUser(context context.Context, token *oauth2.Token) (*User, error) {
	const endpointProfile = "https://www.googleapis.com/oauth2/v2/userinfo"

	response, err := gg.config.Client(context, token).Get(endpointProfile + "?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google responded with a %d trying to fetch user information", response.StatusCode)
	}

	u := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}

	if err = json.NewDecoder(response.Body).Decode(&u); err != nil {
		return nil, err
	}

	user := &User{
		Name:     u.Name,
		Email:    u.Email,
		Image:    u.Picture,
		Username: u.Name, // will be formated by the users service
	}

	return user, nil
}

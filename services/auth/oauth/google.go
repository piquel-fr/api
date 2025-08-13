package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/piquel-fr/api/utils"
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

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	u := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}

	if err := json.Unmarshal(responseBytes, &u); err != nil {
		return nil, err
	}

	user := &User{
		Name:     u.Name,
		Email:    u.Email,
		Image:    u.Picture,
		Username: utils.FormatUsername(u.Name),
	}

	return user, nil
}

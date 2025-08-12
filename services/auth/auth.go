package auth

import "github.com/piquel-fr/api/services/auth/oauth"

func InitAuthService() {
	oauth.InitOAuth()
}

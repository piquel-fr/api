package utils

import (
	"fmt"
	"strings"
	"time"
)

func GenerateSetCookie(name, value, domain, path, sameSite string, age time.Duration) string {
	sameSite = strings.ToLower(sameSite)
	if sameSite != "strict" && sameSite != "lax" {
		sameSite = "none"
	}
	return fmt.Sprintf("%s=%s; Domain=%s; Path=%s; Max-Age=%d; SameSite=%s; Secure; HttpOnly", name, value, domain, path, int(age.Seconds()), sameSite)
}

func GetCookiesFromStr(cookieStr string) map[string]string {
	if cookieStr == "" {
		return nil
	}

	pairs := strings.Split(cookieStr, "; ")

	cookies := map[string]string{}
	for _, pair := range pairs {
		split := strings.SplitN(pair, "=", 2)
		cookies[split[0]] = split[1]
	}

	return cookies
}

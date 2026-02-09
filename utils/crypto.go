package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}


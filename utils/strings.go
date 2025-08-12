package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
	"unicode"
)

func RandString(size int) string {
	nonceBytes := make([]byte, size)
	io.ReadFull(rand.Reader, nonceBytes)
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func HasOnlyLettersAndNumbers(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

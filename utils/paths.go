package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/piquel-fr/api/errors"
)

func IsDir(path string) bool {
	if file, err := os.Stat(path); err != nil {
		return false
	} else {
		return file.IsDir()
	}
}

func ValidatePath(path string) error {
	if strings.Contains(path, "..") {
		return errors.NewError(fmt.Sprintf("Path \"%s\" is not valid as it contains \"..\"", path), http.StatusBadRequest)
	} else if strings.Contains(path, "~") {
		return errors.NewError(fmt.Sprintf("Path \"%s\" is not valid as it contains \"~\"", path), http.StatusBadRequest)
	}
	return nil
}

func FormatLocalPathString(path string) string {
	return fmt.Sprintf("/%s", strings.Trim(path, "/"))
}

func FormatLocalPath(path []byte) []byte {
	return fmt.Appendf(nil, "/%s", bytes.Trim(path, "/"))
}

func HasOnlyLettersAndNumbers(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

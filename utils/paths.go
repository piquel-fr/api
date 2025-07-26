package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

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

func FormatLocalPathString(path, ext string) string {
	trim := strings.Trim(path, "/")
	trim = strings.TrimSuffix(trim, ext)
	return fmt.Sprintf("/%s", trim)
}

func FormatLocalPath(path []byte, ext string) []byte {
	trim := bytes.Trim(path, "/")
	trim = bytes.TrimSuffix(trim, []byte(ext))
	return fmt.Appendf(nil, "/%s", trim)
}

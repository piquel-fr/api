package utils

import "strings"

func FormatUsername(username string) string {
    return strings.ReplaceAll(strings.ToLower(username), " ", "")
}

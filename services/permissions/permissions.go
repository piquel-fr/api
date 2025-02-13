package permissions

import (
	"context"
	"strings"

	"github.com/PiquelChips/piquel.fr/services/database"
)

func UserHasPermission(userId int32, permission string) bool {
    userPermissions, _ := database.Queries.GetUserPermissions(context.Background(), userId)

    for _, userPermission := range userPermissions {
        if matchPermission(permission, userPermission) {
            return true
        }
    }

    return false
}

func matchPermission(permission string, match string) bool {
    permissionParts := strings.Split(permission, ".")
    matchParts := strings.Split(match, ".")

    if len(matchParts) > len(permissionParts) {
        return false
    }

    for i, part := range matchParts {
        if part == "*" {
            return true
        }
        if part != permissionParts[i] {
            return false
        }
    }

    return len(matchParts) == len(permissionParts)
}

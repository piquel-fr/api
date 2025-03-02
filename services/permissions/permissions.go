package permissions

import "github.com/PiquelChips/piquel.fr/utils"

func Authorize(request *Request) error {
	if request.User == nil || request.Ressource == nil {
		// Handle request malformed error
        return nil
	}

	role := request.User.Role
	resourceName := request.Ressource.GetRessourceName()

	if role == "" || resourceName == "" {
		// Handle request malformed error
        return nil
	}

	isAuthozized, err := authorize(request, role, resourceName, []string{})
	if err != nil {
		return err
	}

	if isAuthozized {
		return nil
	}

	// Return access denied error
	return nil
}

func authorize(request *Request, roleName, resourceName string, checkedRoles []string) (bool, error) {
	role, ok := Policy.Roles[roleName]
	if !ok {
		// Handle role not found error
		return false, nil
	}

	var permissions []*Permission

	if role.Permissions == nil {
		permissions = []*Permission{}
	} else {
		permissions = role.Permissions[resourceName]
	}

	parents := role.Parents

	for _, action := range request.Actions {
		if action == "" {
			// Handle request malformed error
			return false, nil
		}

		isAuthozized, err := validateAction(permissions, action, request)
		if err != nil {
			return false, err
		}

		if !isAuthozized && len(parents) > 0 {
			checkedRoles = append(checkedRoles, roleName)

			for _, parent := range parents {
				parentRequest := &Request{
					User:      request.User,
					Ressource: request.Ressource,
					Actions:   []string{action},
				}

				if utils.StringSliceContains(checkedRoles, parent) {
					// Handle role inheritance cycle error
					return false, nil
				}

				isAuthozized, err := authorize(parentRequest, parent, resourceName, checkedRoles)
				if err != nil {
					return false, err
				}

				if isAuthozized {
					break
				}
			}
		}

		if !isAuthozized {
			return false, nil
		}
	}

	return true, nil
}

func validateAction(permissions []*Permission, action string, request *Request) (bool, error) {
	for _, permission := range permissions {
		if permission.Action != action {
			continue
		}

		if len(permission.Conditions) == 0 {
			return true, nil
		}

		isAuthozized, err := checkPermission(permission, request)
		if err != nil {
			return false, err
		}

		if isAuthozized {
			return true, nil
		}
	}

	return false, nil
}

func checkPermission(permission *Permission, request *Request) (bool, error) {
	if permission.Conditions == nil {
		return true, nil
	}

	for _, condition := range permission.Conditions {
		err := condition(request)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

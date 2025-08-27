package auth

import (
	"slices"

	"github.com/piquel-fr/api/utils/errors"
)

func (s *realAuthService) Authorize(request *Request) error {
	if request.User == nil || request.Ressource == nil {
		return newRequestMalformedError(request)
	}

	role := request.User.Role
	resourceName := request.Ressource.GetResourceName()

	if role == "" || resourceName == "" {
		return newRequestMalformedError(request)
	}

	isAuthozized, err := s.authorize(request, role, resourceName, []string{})
	if err != nil {
		return err
	}

	if isAuthozized {
		return nil
	}

	return errors.ErrorForbidden
}

func (s *realAuthService) authorize(request *Request, roleName, resourceName string, checkedRoles []string) (bool, error) {
	role, ok := s.policy.Roles[roleName]
	if !ok {
		return false, newRoleNotFoundError(roleName)
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
			return false, newRequestMalformedError(request)
		}

		isAuthozized, err := s.validateAction(permissions, action, request)
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
					Context:   request.Context,
				}

				if slices.Contains(checkedRoles, parent) {
					return false, newRoleInheritanceCycleError(checkedRoles, parent)
				}

				isAuthozized, err := s.authorize(parentRequest, parent, resourceName, checkedRoles)
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

func (s *realAuthService) validateAction(permissions []*Permission, action string, request *Request) (bool, error) {
	for _, permission := range permissions {

		if permission.Preset != "" {
			permission = s.policy.Permissions[permission.Preset]
		}

		if permission.Action != action {
			continue
		}

		if len(permission.Conditions) == 0 {
			return true, nil
		}

		isAuthozized, err := s.checkPermission(permission, request)
		if err != nil {
			return false, err
		}

		if isAuthozized {
			return true, nil
		}
	}

	return false, nil
}

func (s *realAuthService) checkPermission(permission *Permission, request *Request) (bool, error) {
	if permission.Conditions == nil {
		return true, nil
	}

	// all conditions must pass
	for _, condition := range permission.Conditions {
		err := condition(request)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

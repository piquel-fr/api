package permissions

import (
	"errors"
	"fmt"
)

func newRequestMalformedError(request *Request) error {
	return errors.New(fmt.Sprintf("Request is malformed: %v", request))
}

func newAccessDeniedError() error {
	return errors.New("Access denied!")
}

func newRoleNotFoundError(role string) error {
	return errors.New(fmt.Sprintf("Role %s does not exist!", role))
}

func newRoleInheritanceCycleError(checkedRoles []string, role string) error {
	return errors.New(fmt.Sprintf("There is a role inheritance cycle. Role %s has already been checks: %v.", role, checkedRoles))
}

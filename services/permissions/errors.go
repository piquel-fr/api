package permissions

import (
	"fmt"
)

func newRequestMalformedError(request *Request) error {
	return fmt.Errorf("Request is malformed: %v", request)
}

func newRoleNotFoundError(role string) error {
	return fmt.Errorf("Role %s does not exist!", role)
}

func newRoleInheritanceCycleError(checkedRoles []string, role string) error {
	return fmt.Errorf("There is a role inheritance cycle. Role %s has already been checks: %v.", role, checkedRoles)
}

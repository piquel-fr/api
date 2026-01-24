package auth

import (
	"slices"

	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/email"
	"github.com/piquel-fr/api/utils/errors"
)

const (
	RoleSystem    string = "system"
	RoleAdmin     string = "admin"
	RoleDeveloper string = "developer"
	RoleDefault   string = "default"
)

const (
	ActionView   string = "view"
	ActionCreate string = "create"
	ActionUpdate string = "update"
	ActionDelete string = "delete"
	ActionShare  string = "share"

	ActionListEmailAccounts string = "list_email_accounts"
)

func own(request *Request) error {
	if request.Ressource.GetOwner() == request.User.ID {
		return nil
	}
	return errors.ErrorForbidden
}

func makeOwn(action string) *Permission {
	return &Permission{
		Action:     action,
		Conditions: Conditions{own},
	}
}

var Policy = PolicyConfiguration{
	Permissions: map[string]*Permission{},
	Roles: Roles{
		RoleSystem: {
			Name:        "System",
			Color:       "gray",
			Permissions: map[string][]*Permission{},
			Parents:     []string{RoleDefault, RoleDeveloper, RoleAdmin},
		},
		RoleAdmin: {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*Permission{
				repository.ResourceUser: {
					{Action: ActionUpdate},
					{Action: ActionDelete},
				},
				repository.ResourceMailAccount: {
					{Action: ActionView},
					{Action: ActionUpdate},
					{Action: ActionDelete},
					{Action: ActionListEmailAccounts},
					{Action: ActionShare},
				},
			},
			Parents: []string{RoleDefault, RoleDeveloper},
		},
		RoleDeveloper: {
			Name:  "Developer",
			Color: "blue",
			Permissions: map[string][]*Permission{
				repository.ResourceMailAccount: {
					{
						Action: ActionView,
						Conditions: Conditions{
							func(request *Request) error {
								if request.Ressource.GetOwner() == request.User.ID {
									return nil
								}

								info, ok := request.Ressource.(*email.AccountInfo)
								if !ok {
									return newRequestMalformedError(request)
								}

								if slices.Contains(info.Shares, request.User.Username) {
									return nil
								}
								return errors.ErrorNotFound
							},
						},
					},
					makeOwn(ActionDelete),
				},
				repository.ResourceUser: {
					makeOwn(ActionShare),
					makeOwn(ActionListEmailAccounts),
				},
			},
			Parents: []string{RoleDefault},
		},
		RoleDefault: {
			Name:  "",
			Color: "gray",
			Permissions: map[string][]*Permission{
				repository.ResourceUser: {
					makeOwn(ActionUpdate),
					makeOwn(ActionDelete),
				},
			},
		},
	},
}

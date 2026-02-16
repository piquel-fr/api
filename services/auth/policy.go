package auth

import (
	"slices"

	"github.com/piquel-fr/api/config"
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
	// global actions
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionShare  = "share"

	// admin stuff
	ActionUpdateAdmin = "update_admin"

	// users
	ActionViewEmail = "view_email"

	// sessions
	ActionViewUserSessions   = "view_user_sessions"
	ActionDeleteUserSessions = "delete_user_sessions"

	// email
	ActionListEmailAccounts = "list_email_accounts"
	ActionSendEmail         = "send_email"
)

func own(request *config.AuthRequest) error {
	if request.Ressource.GetOwner() == request.User.ID {
		return nil
	}
	return errors.ErrorForbidden
}

func makeOwn(action string) *config.Permission {
	return &config.Permission{
		Action:     action,
		Conditions: config.Conditions{own},
	}
}

// will check if you own of if the email account is shared with you
func makeOwnEmail(action string) *config.Permission {
	return &config.Permission{
		Action: action,
		Conditions: config.Conditions{
			func(request *config.AuthRequest) error {
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
	}
}

var policy = config.PolicyConfiguration{
	Presets: map[string]*config.Permission{},
	Roles: map[string]*config.Role{
		RoleSystem: {
			Name:        "System",
			Color:       "gray",
			Permissions: map[string][]*config.Permission{},
			Parents:     []string{RoleDefault, RoleDeveloper, RoleAdmin},
		},
		RoleAdmin: {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*config.Permission{
				repository.ResourceUser: {
					{Action: ActionUpdate},
					{Action: ActionDelete},
					{Action: ActionViewEmail},
					{Action: ActionUpdateAdmin},
					{Action: ActionViewUserSessions},
					{Action: ActionDeleteUserSessions},
					{Action: ActionListEmailAccounts},
				},
				repository.ResourceMailAccount: {
					{Action: ActionView},
					{Action: ActionUpdate},
					{Action: ActionDelete},
					{Action: ActionShare},
					{Action: ActionSendEmail},
				},
			},
			Parents: []string{RoleDefault, RoleDeveloper},
		},
		RoleDeveloper: {
			Name:  "Developer",
			Color: "blue",
			Permissions: map[string][]*config.Permission{
				repository.ResourceUser: {
					makeOwn(ActionListEmailAccounts),
				},
				repository.ResourceMailAccount: {
					makeOwnEmail(ActionView),
					makeOwnEmail(ActionSendEmail),
					makeOwn(ActionShare),
					makeOwn(ActionDelete),
				},
			},
			Parents: []string{RoleDefault},
		},
		RoleDefault: {
			Name:  "",
			Color: "gray",
			Permissions: map[string][]*config.Permission{
				repository.ResourceUser: {
					makeOwn(ActionUpdate),
					makeOwn(ActionDelete),
					makeOwn(ActionViewEmail),
					makeOwn(ActionViewUserSessions),
					makeOwn(ActionDeleteUserSessions),
				},
			},
		},
	},
}

package auth

import "github.com/PiquelChips/piquel.fr/errors"

var Policy = &PolicyConfiguration{
	Permissions: map[string]*Permission{
		"updateOwn": {
			Action: "update",
			Conditions: Conditions{
				func(request *Request) error {
					if request.Ressource.GetOwner() == request.User.ID {
						return nil
					}
					return errors.ErrorNotAuthenticated
				},
			},
		},
	},
	Roles: Roles{
		"admin": {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*Permission{
				"user": {
					{Action: "create"},
					{Action: "update"},
					{Action: "delete"},
				},
			},
			Parents: []string{"default", "developer"},
		},
		"developer": {
			Name:    "Developer",
			Color:   "blue",
			Parents: []string{"default"},
		},
		"default": {
			Name:  "",
			Color: "gray",
			Permissions: map[string][]*Permission{
				"user": {
					{Preset: "updateOwn"},
				},
			},
		},
	},
}

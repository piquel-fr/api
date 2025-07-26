package auth

import (
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
)

var Policy = &PolicyConfiguration{
	Permissions: map[string]*Permission{
		"updateOwn": {
			Action: "update",
			Conditions: Conditions{
				func(request *Request) error {
					if request.Ressource.GetOwner() == request.User.ID {
						return nil
					}
					return errors.ErrorForbidden
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
				"documentation": {
					{Action: "view"},
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
				"documentation": {
					{
						Action: "view",
						Conditions: Conditions{
							func(request *Request) error {
								docs, ok := request.Ressource.(*models.Documentation)
								if !ok {
									return newRequestMalformedError(request)
								}

								if docs.Public {
									return nil
								} else if docs.GetOwner() == request.User.ID {
									return nil
								}

								return errors.ErrorForbidden
							},
						},
					},
				},
			},
		},
	},
}

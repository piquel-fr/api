package auth

import (
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/database"
)

func own(request *Request) error {
	if request.Ressource.GetOwner() == request.User.ID {
		return nil
	}
	return errors.ErrorForbidden
}

var Policy = &PolicyConfiguration{
	Permissions: map[string]*Permission{
		"updateOwn": {
			Action:     "update",
			Conditions: Conditions{own},
		},
		"deleteOwn": {
			Action:     "delete",
			Conditions: Conditions{own},
		},
	},
	Roles: Roles{
		"admin": {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*Permission{
				"user": {
					{Action: "update"},
					{Action: "delete"},
				},
				"documentation": {
					{Action: "view"},
					{Action: "create"},
					{Action: "update"},
					{Action: "delete"},
				},
			},
			Parents: []string{"default", "developer"},
		},
		"default": {
			Name:  "",
			Color: "gray",
			Permissions: map[string][]*Permission{
				"user": {
					{Preset: "updateOwn"},
					{Preset: "deleteOwn"},
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
					{Preset: "updateOwn"},
					{Preset: "deleteOwn"},
					{
						Action: "create",
						Conditions: Conditions{
							func(request *Request) error {
								count, err := database.Queries.CountUserDocumentations(request.Context, request.Ressource.GetOwner())
								if err != nil {
									return err
								}

								if count >= config.Configuration.MaxDocumentationCount {
									return errors.NewError(
										fmt.Sprintf("you already have %d/%d documentation instances", count, config.Configuration.MaxDocumentationCount),
										http.StatusForbidden,
									)
								}

								return nil
							},
						},
					},
				},
			},
		},
	},
}

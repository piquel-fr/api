package auth

import (
	"fmt"
	"net/http"

	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/utils/errors"
)

func own(request *Request) error {
	if request.Ressource.GetOwner() == request.User.ID {
		return nil
	}
	return errors.ErrorForbidden
}

func (s *realAuthService) createPolicy() {
	s.policy = PolicyConfiguration{
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
					"docs_instance": {
						{
							Action: "view",
							Conditions: Conditions{
								func(request *Request) error {
									docs, ok := request.Ressource.(*models.DocsInstance)
									if !ok {
										return newRequestMalformedError(request)
									}

									if docs.Public {
										return nil
									}

									if docs.GetOwner() == request.User.ID {
										return nil
									}

									return errors.ErrorForbidden
								},
							},
						},
						{
							Action: "create",
							Conditions: Conditions{
								func(request *Request) error {
									count, err := s.database.CountUserDocsInstances(request.Context, request.User.ID)
									if err != nil {
										return err
									}

									if count >= s.config.MaxDocsInstanceCount {
										return errors.NewError(
											fmt.Sprintf("you already have %d/%d documentation instances", count, s.config.MaxDocsInstanceCount),
											http.StatusForbidden,
										)
									}

									return nil
								},
							},
						},
						{Preset: "updateOwn"},
						{Preset: "deleteOwn"},
					},
				},
			},
		},
	}
}

package auth

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/utils/errors"
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

const RoleSystem string = "system"

var policy = PolicyConfiguration{
	Permissions: map[string]*Permission{},
	Roles: Roles{
		RoleSystem: {
			Name:        "System",
			Color:       "gray",
			Permissions: map[string][]*Permission{},
			Parents:     []string{"default", "developer", "admin"},
		},
		"admin": {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*Permission{
				repository.ResourceUser: {
					{Action: "update"},
					{Action: "delete"},
				},
				repository.ResourceDocsInstance: {
					{Action: "view"},
					{Action: "create"},
					{Action: "update"},
					{Action: "delete"},
				},
				repository.ResourceMailAccount: {
					{Action: "view"},
					{Action: "update"},
					{Action: "delete"},
					{Action: "list_email_accounts"},
					{Action: "share"},
				},
			},
			Parents: []string{"default", "developer"},
		},
		"developer": {
			Name:  "Developer",
			Color: "blue",
			Permissions: map[string][]*Permission{
				repository.ResourceMailAccount: {
					{
						Action: "view",
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
					makeOwn("delete"),
				},
				repository.ResourceUser: {
					makeOwn("share"),
					makeOwn("list_email_accounts"),
				},
			},
			Parents: []string{"default"},
		},
		"default": {
			Name:  "",
			Color: "gray",
			Permissions: map[string][]*Permission{
				repository.ResourceUser: {
					makeOwn("update"),
					makeOwn("delete"),
				},
				repository.ResourceDocsInstance: {
					{
						Action: "view",
						Conditions: Conditions{
							func(request *Request) error {
								docs, ok := request.Ressource.(*repository.DocsInstance)
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
								count, err := database.Queries.CountUserDocsInstances(request.Context, request.User.ID)
								if err != nil {
									return err
								}

								if count >= config.MaxDocsInstanceCount {
									return errors.NewError(
										fmt.Sprintf("you already have %d/%d documentation instances", count, config.MaxDocsInstanceCount),
										http.StatusForbidden,
									)
								}

								return nil
							},
						},
					},
					makeOwn("update"),
					makeOwn("delete"),
				},
			},
		},
	},
}

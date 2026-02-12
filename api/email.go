package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/email"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type EmailHandler struct {
	userService  users.UserService
	authService  auth.AuthService
	emailService email.EmailService
}

func CreateEmailHandler(userService users.UserService, authService auth.AuthService, emailService email.EmailService) *EmailHandler {
	return &EmailHandler{userService, authService, emailService}
}

func (h *EmailHandler) getName() string { return "email" }

func (h *EmailHandler) getSpec() Spec {
	spec := newSpecBase(h)

	accountSchema := openapi3.NewObjectSchema().
		WithProperty("id", openapi3.NewInt32Schema()).
		WithProperty("ownerId", openapi3.NewInt32Schema()).
		WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("name", openapi3.NewStringSchema())

	addAccountSchema := openapi3.NewObjectSchema().
		WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("password", openapi3.NewStringSchema())

	// TODO: update
	folderSchema := openapi3.NewObjectSchema().
		WithProperty("name", openapi3.NewStringSchema())

	// TODO: update
	emailMessageSchema := openapi3.NewObjectSchema().
		WithProperty("id", openapi3.NewStringSchema()).
		WithProperty("subject", openapi3.NewStringSchema()).
		WithProperty("from", openapi3.NewStringSchema()).
		WithProperty("body", openapi3.NewStringSchema())

	// TODO: update
	sendEmailPayloadSchema := openapi3.NewObjectSchema().
		WithProperty("to", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("subject", openapi3.NewStringSchema()).
		WithProperty("body", openapi3.NewStringSchema())

	spec.Components.Schemas = openapi3.Schemas{
		"MailAccount":       &openapi3.SchemaRef{Value: accountSchema},
		"AddAccountPayload": &openapi3.SchemaRef{Value: addAccountSchema},
		"Folder":            &openapi3.SchemaRef{Value: folderSchema},
		"EmailMessage":      &openapi3.SchemaRef{Value: emailMessageSchema},
		"SendEmailPayload":  &openapi3.SchemaRef{Value: sendEmailPayloadSchema},
	}

	spec.AddOperation("/", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "List email accounts",
		Description: "List accounts belonging to the authenticated user, or a specific user if admin.",
		OperationID: "list-email-accounts",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "query",
					Required:    false,
					Description: "Optional username to filter by (admin only)",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "count",
					In:          "query",
					Required:    false,
					Description: "If present, returns only the count of accounts as text",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewBoolSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("List of email accounts").
					WithJSONSchemaRef(openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(accountSchema))),
			}),
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewInt32Schema(), nil)).WithDescription("Count of all the email accounts"),
			}),
		),
	})

	spec.AddOperation("/", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Create email account",
		Description: "Create a new email account for the authenticated user",
		OperationID: "add-email-account",
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: true,
				Content: openapi3.NewContentWithJSONSchemaRef(
					openapi3.NewSchemaRef("#/components/schemas/AddAccountPayload", addAccountSchema),
				),
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Account created successfully")}),
			openapi3.WithStatus(400, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewStringSchema(), []string{"text/plain"})).WithDescription("Invalid input")}),
		),
	})

	spec.AddOperation("/{email}", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Get account info",
		Description: "Get details of a specific email account",
		OperationID: "get-email-account",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "email",
					In:          "path",
					Required:    true,
					Description: "The email address of the account",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		// TODO: fix this returns account info not just an account
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("Account details").
					WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/MailAccount", accountSchema)),
			}),
		),
	})

	spec.AddOperation("/{email}", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Send email",
		Description: "Send an email from this account",
		OperationID: "send-email",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: true,
				Content: openapi3.NewContentWithJSONSchemaRef(
					openapi3.NewSchemaRef("#/components/schemas/SendEmailPayload", sendEmailPayloadSchema),
				),
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Email sent successfully")}),
		),
	})

	spec.AddOperation("/{email}", http.MethodDelete, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Remove account",
		Description: "Delete an email account",
		OperationID: "delete-email-account",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "email",
					In:          "path",
					Required:    true,
					Description: "The email address to delete",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Account deleted successfully")}),
		),
	})

	spec.AddOperation("/{email}/share", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Share account",
		Description: "Share an email account with another user",
		OperationID: "share-email-account",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "query",
					Required:    true,
					Description: "The username to share the account with",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Account shared successfully")}),
		),
	})

	spec.AddOperation("/{email}/share", http.MethodDelete, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Remove share",
		Description: "Stop sharing an email account with a specific user",
		OperationID: "unshare-email-account",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "query",
					Required:    true,
					Description: "The username to remove access from",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Share removed successfully")}),
		),
	})

	spec.AddOperation("/{email}/folder", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "List folders",
		Description: "Get a list of folders for the email account",
		OperationID: "list-folders",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("List of folders").
					WithJSONSchemaRef(openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(folderSchema))),
			}),
		),
	})

	spec.AddOperation("/{email}/folder", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Create folder",
		Description: "Create a new folder",
		OperationID: "create-folder",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "name",
					In:          "query",
					Required:    true,
					Description: "The name of the folder to create",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Folder created successfully")}),
		),
	})

	spec.AddOperation("/{email}/folder/{folder}", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "List emails",
		Description: "List emails within a specific folder with pagination",
		OperationID: "list-folder-emails",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "folder",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "offset",
					In:          "query",
					Required:    false,
					Description: "The number of items to skip before starting to collect the result set",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewInt32Schema().WithMin(0)},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "limit",
					In:          "query",
					Required:    false,
					Description: "The numbers of items to return (max 200)",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewInt32Schema().WithMin(1).WithMax(200)},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("List of emails").
					WithJSONSchemaRef(openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(emailMessageSchema))),
			}),
		),
	})

	spec.AddOperation("/{email}/folder/{folder}", http.MethodDelete, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Delete folder",
		Description: "Delete a folder",
		OperationID: "delete-folder",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "folder",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Folder deleted successfully")}),
		),
	})

	spec.AddOperation("/{email}/folder/{folder}", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Rename folder",
		Description: "Rename a specific folder",
		OperationID: "rename-folder",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "folder",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "name",
					In:          "query",
					Required:    true,
					Description: "The new name for the folder",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Folder renamed successfully")}),
		),
	})

	spec.AddOperation("/{email}/folder/{folder}/{id}", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Get email",
		Description: "Get a specific email message",
		OperationID: "get-email",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "folder",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewInt32Schema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("Email details").
					WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/EmailMessage", emailMessageSchema)),
			}),
		),
	})

	spec.AddOperation("/{email}/folder/{folder}/{id}", http.MethodDelete, &openapi3.Operation{
		Tags:        []string{"email"},
		Summary:     "Delete email",
		Description: "Delete a specific email message",
		OperationID: "delete-email",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "email",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "folder",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("Email deleted successfully")}),
		),
	})

	return spec
}

func (h *EmailHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("PUT /", h.handleAddAccount)
	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET", "PUT"))

	handler.HandleFunc("GET /{email}", h.handleAccountInfo)
	handler.HandleFunc("PUT /{email}", h.handleSendEmail)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)
	handler.Handle("OPTIONS /{email}", middleware.CreateOptionsHandler("GET", "PUT", "DELETE"))

	// sharing
	handler.HandleFunc("PUT /{email}/share", h.handleShareAccount)
	handler.HandleFunc("DELETE /{email}/share", h.handleRemoveAccountShare)
	handler.Handle("OPTIONS /{email}/share", middleware.CreateOptionsHandler("PUT", "DELETE"))

	handler.HandleFunc("GET /{email}/folder", h.handleGetFolders)
	handler.HandleFunc("PUT /{email}/folder", h.handleCreateFolder)
	handler.Handle("OPTIONS /{email}/folder", middleware.CreateOptionsHandler("GET", "PUT"))

	handler.HandleFunc("GET /{email}/folder/{folder}", h.handleGetFolderEmails)
	handler.HandleFunc("DELETE /{email}/folder/{folder}", h.handleDeleteFolder)
	handler.HandleFunc("PUT /{email}/folder/{folder}", h.handleRenameFolder)
	handler.Handle("OPTIONS /{email}/folder/{folder}", middleware.CreateOptionsHandler("GET", "DELETE", "PUT"))

	handler.HandleFunc("GET /{email}/folder/{folder}/{id}", h.handleGetEmail)
	handler.HandleFunc("DELETE /{email}/folder/{folder}/{id}", h.handleDeleteEmail)
	handler.Handle("OPTIONS /{email}/folder/{folder}/{id}", middleware.CreateOptionsHandler("GET", "DELETE"))

	// emails
	return handler
}

func (h *EmailHandler) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	requester, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	var user *repository.User
	if username := r.URL.Query().Get("user"); username != "" {
		user, err = h.userService.GetUserByUsername(r.Context(), username)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
	} else {
		user = requester
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      requester,
		Ressource: user,
		Context:   r.Context(),
		Actions:   []string{auth.ActionListEmailAccounts},
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("count") {
		count, err := h.emailService.CountAccounts(r.Context(), user.ID)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(int(count))))
		return
	}

	accounts, err := h.emailService.ListAccounts(r.Context(), user.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	for i := range accounts {
		accounts[i].Username = ""
		accounts[i].Password = ""
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleAddAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.AddEmailAccountParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.OwnerId = user.ID
	if _, err = database.Queries.AddEmailAccount(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleAccountInfo(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	accountInfo, err := h.emailService.GetAccountInfo(r.Context(), account)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: &accountInfo,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	accountInfo.Username = ""
	accountInfo.Password = ""

	data, err := json.Marshal(accountInfo)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionSendEmail},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := email.EmailSendParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.SendEmail(account, params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionDelete},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.RemoveAccount(r.Context(), account.ID); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleShareAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionShare},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	sharingUser, err := h.userService.GetUserByUsername(r.Context(), r.URL.Query().Get("user"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.AddShareParams{
		UserId:     sharingUser.ID,
		Account:    account.ID,
		Permission: "",
	}

	if err := h.emailService.AddShare(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleRemoveAccountShare(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionShare},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	sharingUser, err := h.userService.GetUserByUsername(r.Context(), r.URL.Query().Get("user"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.RemoveShare(r.Context(), sharingUser.ID, account.ID); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleGetFolders(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	folders, err := h.emailService.ListFolders(account)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(folders)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.CreateFolder(account, r.URL.Query().Get("name")); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) handleGetFolderEmails(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	limit, offset := config.MaxLimit, 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = config.MaxLimit
		}
		if limit < 1 {
			limit = 1
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}
	}

	emails, err := h.emailService.GetFolderEmails(account, r.PathValue("folder"), uint32(offset), uint32(limit))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(emails)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleDeleteFolder(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.DeleteFolder(account, r.PathValue("folder")); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) handleRenameFolder(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.RenameFolder(account, r.PathValue("folder"), r.URL.Query().Get("name")); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *EmailHandler) handleGetEmail(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	email, err := h.emailService.GetEmail(account, r.PathValue("folder"), uint32(id))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(email)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleDeleteEmail(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&config.AuthRequest{
		User:      user,
		Ressource: account,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.DeleteEmail(account, r.PathValue("folder"), uint32(id)); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

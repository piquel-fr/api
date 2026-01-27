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

	spec.Components.Schemas = openapi3.Schemas{
		"MailAccount":       &openapi3.SchemaRef{Value: accountSchema},
		"AddAccountPayload": &openapi3.SchemaRef{Value: addAccountSchema},
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
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("Account details").
					WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/MailAccount", accountSchema)),
			}),
		),
	})

	// 6. Operation: Remove Account (DELETE /{email})
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

	// 7. Operation: Share Account (PUT /{email}/share)
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

	// 8. Operation: Remove Share (DELETE /{email}/share)
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

	return spec
}

func (h *EmailHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("PUT /", h.handleAddAccount)
	handler.HandleFunc("GET /{email}", h.handleAccountInfo)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)

	// sharing
	handler.HandleFunc("PUT /{email}/share", h.handleShareAccount)
	handler.HandleFunc("DELETE /{email}/share", h.handleRemoveAccountShare)

	// OPTIONS handlers
	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET", "PUT"))
	handler.Handle("OPTIONS /{email}", middleware.CreateOptionsHandler("GET", "DELETE"))
	handler.Handle("OPTIONS /{email}/share", middleware.CreateOptionsHandler("PUT", "DELETE"))

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

	for _, account := range accounts {
		account.Username = ""
		account.Password = ""
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

	accountInfo, err := h.emailService.GetAccountInfo(r.Context(), &account)
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
		Ressource: &account,
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
		Ressource: &account,
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
		Ressource: &account,
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

	params := repository.RemoveShareParams{
		UserId:  sharingUser.ID,
		Account: account.ID,
	}

	if err := h.emailService.RemoveShare(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

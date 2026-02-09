package api

import (
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type UserHandler struct {
	userService users.UserService
	authService auth.AuthService
}

func CreateUserHandler(userService users.UserService, authService auth.AuthService) *UserHandler {
	return &UserHandler{userService, authService}
}

func (h *UserHandler) getName() string { return "users" }

func (h *UserHandler) getSpec() Spec {
	spec := newSpecBase(h)

	userSchema := openapi3.NewObjectSchema().
		WithProperty("id", openapi3.NewInt32Schema()).
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("image", openapi3.NewStringSchema()).
		WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("role", openapi3.NewStringSchema()).
		WithProperty("createdAt", openapi3.NewDateTimeSchema()).
		WithRequired([]string{"id", "username", "name", "image", "role", "createdAt"})

	updateUserSchema := openapi3.NewObjectSchema().
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("image", openapi3.NewStringSchema()).
		WithRequired([]string{"username", "name", "image"})

	updateUserAdminSchema := openapi3.NewObjectSchema().
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("image", openapi3.NewStringSchema()).
		WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("role", openapi3.NewStringSchema()).
		WithRequired([]string{"username", "name", "image", "email", "role"})

	spec.Components.Schemas = openapi3.Schemas{
		"User":                  &openapi3.SchemaRef{Value: userSchema},
		"UpdateUserParams":      &openapi3.SchemaRef{Value: updateUserSchema},
		"UpdateUserAdminParams": &openapi3.SchemaRef{Value: updateUserAdminSchema},
	}

	spec.AddOperation("/self", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"users"},
		Summary:     "Get self user object",
		Description: "Get the user that is associated with the auth",
		OperationID: "get-self",
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().WithDescription("User profile found").WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/User", userSchema)),
			}),
		),
	})

	spec.AddOperation("/{user}", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"users"},
		Summary:     "Get specific user",
		Description: "Get the profile of the user specified in the path",
		OperationID: "get-user-by-path",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "path",
					Required:    true,
					Description: "The username",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{
				Value: openapi3.NewResponse().WithDescription("User profile found").WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/User", userSchema)),
			}),
		),
	})

	spec.AddOperation("/{user}", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"users"},
		Summary:     "Update user",
		Description: "Update the profile of the specified user",
		OperationID: "update-user",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "path",
					Required:    true,
					Description: "The username to update",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: true,
				Content: openapi3.NewContentWithJSONSchemaRef(
					openapi3.NewSchemaRef("#/components/schemas/UpdateUserParams", updateUserSchema),
				),
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("User updated successfully")}),
			openapi3.WithStatus(400, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewStringSchema(), []string{"text/plain"})).WithDescription("Invalid input or json")}),
		),
	})

	spec.AddOperation("/{user}", http.MethodDelete, &openapi3.Operation{
		Tags:        []string{"users"},
		Summary:     "Delete user",
		Description: "Delete the user specified in the path",
		OperationID: "delete-user",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "path",
					Required:    true,
					Description: "The username",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("User deleted successfully")}),
			openapi3.WithStatus(401, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewStringSchema(), []string{"text/plain"})).WithDescription("Unauthorized")}),
		),
	})

	spec.AddOperation("/{user}/admin", http.MethodPut, &openapi3.Operation{
		Tags:        []string{"users", "admin"},
		Summary:     "Update user",
		Description: "Update the profile of the specified user",
		OperationID: "update-user-admin",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "user",
					In:          "path",
					Required:    true,
					Description: "The username to update",
					Schema:      &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
				},
			},
		},
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: true,
				Content: openapi3.NewContentWithJSONSchemaRef(
					openapi3.NewSchemaRef("#/components/schemas/UpdateUserAdminParams", updateUserAdminSchema),
				),
			},
		},
		Responses: openapi3.NewResponses(
			openapi3.WithStatus(200, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("User updated successfully")}),
			openapi3.WithStatus(400, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewStringSchema(), []string{"text/plain"})).WithDescription("Invalid input or json")}),
		),
	})

	return spec
}

func (h *UserHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /self", h.handleGetSelf)
	handler.Handle("OPTIONS /self", middleware.CreateOptionsHandler("GET"))

	handler.HandleFunc("GET /{user}", h.handleGetUser)
	handler.HandleFunc("PUT /{user}", h.handlePutUser)
	handler.HandleFunc("DELETE /{user}", h.handleDeleteUser)
	handler.Handle("OPTIONS /{user}", middleware.CreateOptionsHandler("GET", "PUT", "DELETE"))

	handler.HandleFunc("PUT /{user}/admin", h.handlePutUserAdmin)
	handler.Handle("OPTIONS /{user}/admin", middleware.CreateOptionsHandler("PUT"))

	// TODO: session management endpoints

	return handler
}

func (h *UserHandler) handleGetSelf(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")
	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	requester, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &config.AuthRequest{
		User:      requester,
		Ressource: user,
		Actions:   []string{auth.ActionViewEmail},
		Context:   r.Context(),
	}

	if err := h.authService.Authorize(request); err == errors.ErrorForbidden {
		user.Email = ""
	} else if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handlePutUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	requester, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &config.AuthRequest{
		User:      requester,
		Ressource: user,
		Actions:   []string{auth.ActionUpdate},
		Context:   r.Context(),
	}

	if err := h.authService.Authorize(request); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.UpdateUserParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.ID = user.ID

	if err := h.userService.UpdateUser(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")
	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	requester, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &config.AuthRequest{
		User:      requester,
		Ressource: user,
		Actions:   []string{auth.ActionDelete},
		Context:   r.Context(),
	}

	if err := h.authService.Authorize(request); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.userService.DeleteUser(r.Context(), user); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) handlePutUserAdmin(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	requester, err := h.userService.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &config.AuthRequest{
		User:      requester,
		Ressource: user,
		Actions:   []string{auth.ActionUpdateAdmin},
		Context:   r.Context(),
	}

	if err := h.authService.Authorize(request); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.UpdateUserAdminParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.ID = user.ID

	if err := h.userService.UpdateUserAdmin(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type ProfileHandler struct {
	userService users.UserService
	authService auth.AuthService
}

func CreateProfileHandler(userService users.UserService, authService auth.AuthService) *ProfileHandler {
	return &ProfileHandler{userService, authService}
}

func (h *ProfileHandler) getName() string { return "profile" }

func (h *ProfileHandler) getSpec() Spec {
	spec := newSpecBase(h)
	spec.Info.Description = "DEPRECATED " + spec.Info.Description

	userSchema := openapi3.NewObjectSchema().
		WithProperty("id", openapi3.NewInt32Schema()).
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("image", openapi3.NewStringSchema()).
		WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
		WithProperty("role", openapi3.NewStringSchema()).
		WithProperty("createdAt", openapi3.NewDateTimeSchema()).
		WithRequired([]string{"id", "username", "name", "image", "email", "role", "createdAt"})

	updateUserSchema := openapi3.NewObjectSchema().
		WithProperty("username", openapi3.NewStringSchema()).
		WithProperty("name", openapi3.NewStringSchema()).
		WithProperty("image", openapi3.NewStringSchema()).
		WithRequired([]string{"username", "name", "image"})

	spec.Components.Schemas = openapi3.Schemas{
		"User":             &openapi3.SchemaRef{Value: userSchema},
		"UpdateUserParams": &openapi3.SchemaRef{Value: updateUserSchema},
	}

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
			openapi3.WithStatus(401, &openapi3.ResponseRef{Value: openapi3.NewResponse().WithContent(openapi3.NewContentWithSchema(openapi3.NewStringSchema(), []string{"text/plain"})).WithDescription("Unauthorized")}),
		),
	})

	spec.AddOperation("/", http.MethodGet, &openapi3.Operation{
		Tags:        []string{"users"},
		Summary:     "Get user profile",
		Description: "Get user by query param 'username', or the currently authenticated user if empty",
		OperationID: "get-profile",
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:        "username",
					In:          "query",
					Required:    false,
					Description: "Optional username. If omitted, returns current user.",
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

	return spec
}

func (h *ProfileHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", h.handleGetProfileQuery)
	handler.HandleFunc("GET /{user}", h.handleGetProfile)
	handler.HandleFunc("PUT /{user}", h.handleUpdateProfile)

	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{user}", middleware.CreateOptionsHandler("GET", "PUT"))

	return handler
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	h.writeProfile(w, r, r.PathValue("user"))
}

func (h *ProfileHandler) handleGetProfileQuery(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		user, err := h.authService.GetUserFromContext(r.Context())
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
		username = user.Username
	}

	h.writeProfile(w, r, username)
}

func (h *ProfileHandler) writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &config.AuthRequest{
		User:      user,
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

	if err := database.Queries.UpdateUser(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

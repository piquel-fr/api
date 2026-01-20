package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/middleware"
)

func GetHumaConfig(title, version string) huma.Config {
	schemaPrefix := "#/components/schemas/"

	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	return huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   title,
				Version: version,
				Contact: &huma.Contact{
					Name:  "API Support",
					URL:   "https://piquel.fr",
					Email: "contact@piquel.fr",
				},
			},
			Components: &huma.Components{
				Schemas: registry,
			},
			Servers: []*huma.Server{
				{
					URL:         fmt.Sprintf("%s/profile", config.Envs.Url),
					Description: "Main production endpoint of the API",
				},
			},
		},
		OpenAPIPath:   "/openapi",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks: []func(huma.Config) huma.Config{
			func(c huma.Config) huma.Config {
				// Add a link transformer to the API. This adds `Link` headers and
				// puts `$schema` fields in the response body which point to the JSON
				// Schema that describes the response structure.
				// This is a create hook so we get the latest schema path setting.
				linkTransformer := huma.NewSchemaLinkTransformer(schemaPrefix, c.SchemasPath)
				c.OnAddOperation = append(c.OnAddOperation, linkTransformer.OnAddOperation)
				c.Transformers = append(c.Transformers, linkTransformer.Transform)
				return c
			},
		},
	}
}

type ProfileHandler struct {
	api         huma.API
	authService auth.AuthService
}

func CreateProfileHandler(router humago.Mux, authService auth.AuthService) ProfileHandler {
	handler := ProfileHandler{
		api:         humago.New(router, GetHumaConfig("Piquel Profile API", "0.1.0")),
		authService: authService,
	}

	handler.api.UseMiddleware(middleware.AuthMiddleware(handler.authService, handler.api))

	huma.Register(handler.api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        "/",
	}, handler.handleGetProfile)
	return handler
}

func (h *ProfileHandler) handleGetProfile(ctx context.Context, in *struct {
	User string `query:"user"`
}) (*struct{}, error) {

	if in.User == "" {
		id, err := h.authService.GetUserIdFromContext(ctx)
		if err != nil {
			return nil, err
		}
		profile, err := h.authService.GetProfileFromUserId(ctx, id)
		if err != nil {
			return nil, err
		}
		in.User = profile.Username
	}

	return nil, nil
}

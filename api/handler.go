package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/email"
	"github.com/piquel-fr/api/utils/middleware"
)

type Handler interface {
	getName() string
	getSpec() *openapi3.T
	createHttpHandler() http.Handler
}

func CreateRouter(authService auth.AuthService, emailService email.EmailService) (http.Handler, error) {
	// these routes are unauthenticated and should remail so.
	// do not any other routes to this router. all other routes
	// should be added to createProtectedRouter
	router := http.NewServeMux()
	router.HandleFunc("/{$}", rootHandler)
	router.Handle("/auth/", http.StripPrefix("/auth", CreateAuthHandler(authService).createHttpHandler()))

	handlers := []Handler{
		CreateProfileHandler(authService),
		CreateEmailHandler(authService, emailService),
	}

	for _, handler := range handlers {
		// spec
		specPath := fmt.Sprintf("/specification/%s.json", handler.getName())
		specHandler, err := newSpecHandler(handler.getSpec())
		if err != nil {
			return nil, err
		}
		router.HandleFunc(specPath, specHandler)
	}

	// bind the protected router
	protectedRouter := createProtectedRouter(handlers)
	protectedRouter = middleware.AddMiddleware(protectedRouter, middleware.AuthMiddleware(authService))
	router.Handle("/", protectedRouter)

	return middleware.AddMiddleware(router, middleware.CORSMiddleware), nil
}

func createProtectedRouter(handlers []Handler) http.Handler {
	router := http.NewServeMux()
	for _, handler := range handlers {
		log.Printf("Registering %s handler", handler.getName())
		// handler
		prefix := fmt.Sprintf("/%s", handler.getName())
		path := fmt.Sprintf("/%s/", handler.getName())
		router.Handle(path, http.StripPrefix(prefix, handler.createHttpHandler()))
	}
	return router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Welcome to the Piquel API! Visit the <a href=\"https://piquel.fr/docs\">API</a> for more information."))
	w.WriteHeader(http.StatusOK)
}

func newSpecHandler(spec *openapi3.T) (http.HandlerFunc, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}, nil
}

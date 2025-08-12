package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/piquel-fr/api/handlers"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/docs"
	gh "github.com/piquel-fr/api/services/github"
	"github.com/piquel-fr/api/services/middleware"
)

func main() {
	log.Printf("Initializing piquel.fr API...\n")

	// Intialize services
	config.LoadConfig()
	gh.InitGithubWrapper()
	auth.InitAuthentication()
	auth.InitCookieStore()
	database.InitDatabase()
	defer database.DeinitDatabase()
	docs.InitDocsService()

	models.Init()

	// Initialize the router
	router := mux.NewRouter()
	middleware.Setup(router)

	log.Printf("[Router] Initialized router!\n")

	router.HandleFunc("/profile", handlers.HandleGetProfileQuery).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/profile/{profile}", handlers.HandleGetProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/profile/{profile}", handlers.HandleUpdateProfile).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/auth/logout", handlers.HandleLogout).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}", handlers.HandleProviderLogin).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}/callback", handlers.HandleAuthCallback).Methods(http.MethodGet, http.MethodOptions)

	// extra methods are for CORS middleware. Will be fixed when moving to net/http
	router.HandleFunc("/docs", handlers.HandleListDocs).Methods(http.MethodGet, http.MethodOptions, http.MethodPost) // GET
	router.HandleFunc("/docs", handlers.HandleNewDocs).Methods(http.MethodPost, http.MethodOptions)                  // POST
	// extra methods are for CORS middleware. Will be fixed when moving to net/http
	router.HandleFunc("/docs/{documentation}", handlers.HandleGetDocs).Methods(http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete) // GET
	router.HandleFunc("/docs/{documentation}", handlers.HandleUpdateDocs).Methods(http.MethodPut, http.MethodOptions)                                 // PUT
	router.HandleFunc("/docs/{documentation}", handlers.HandleDeleteDocs).Methods(http.MethodDelete, http.MethodOptions)                              // DELETE
	router.PathPrefix("/docs/{documentation}/page").HandlerFunc(handlers.HandleGetDocsPage).Methods(http.MethodGet, http.MethodOptions)               // GET

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", http.ListenAndServe(address, router).Error())
}

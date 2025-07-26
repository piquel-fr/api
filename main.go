package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/piquel-fr/api/handlers"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/config"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/docs"
	gh "github.com/piquel-fr/api/services/github"
	"github.com/piquel-fr/api/services/middleware"
	"github.com/piquel-fr/api/models"
)

func main() {
	log.Printf("Initializing piquel.fr API...\n")

	// Intialize services
	config.LoadConfig()
	auth.InitAuthentication()
	auth.InitCookieStore()
	database.InitDatabase()
	defer database.DeinitDatabase()
	gh.InitGithubWrapper()

	if err := docs.InitDocumentation(); err != nil {
		panic(err)
	}

	models.Init()

	// Initialize the router
	router := mux.NewRouter()
	middleware.Setup(router)

	log.Printf("[Router] Initialized router!\n")

	router.HandleFunc("/profile", handlers.HandleBaseProfile).Methods(http.MethodGet, http.MethodPut, http.MethodOptions)
	router.HandleFunc("/profile/{profile}", handlers.HandleProfile).Methods(http.MethodGet, http.MethodPut, http.MethodOptions)

	router.HandleFunc("/auth/logout", handlers.HandleLogout).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}", handlers.HandleProviderLogin).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}/callback", handlers.HandleAuthCallback).Methods(http.MethodGet, http.MethodOptions)

	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.HandlerFunc(handlers.HandleDocs)))

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", http.ListenAndServe(address, router).Error())
}

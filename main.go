package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/piquel-fr/api/handlers"
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
	auth.InitCookieStore()
	database.InitDatabase()
	defer database.DeinitDatabase()
	docs.InitDocsService()
	auth.InitAuthService()

	// Initialize the router
	router := http.NewServeMux()

	log.Printf("[Router] Initialized router!\n")

	router.Handle("/profile/", http.StripPrefix("/profile", handlers.CreateProfileHandler()))
	router.Handle("/auth/", http.StripPrefix("/auth", handlers.CreateAuthHandler()))
	router.Handle("/docs/", http.StripPrefix("/docs", handlers.CreateDocsHandler()))

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)

	server := http.Server{
		Addr: address,
		Handler: middleware.AddMiddleware(router,
			middleware.CORSMiddleware,
			middleware.AuthMiddleware,
		),
	}

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", server.ListenAndServe().Error())
}

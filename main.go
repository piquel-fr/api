package main

import (
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/handlers"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Initializing piquel.fr website...\n")

    // Intialize services
    config.LoadConfig()
    auth.InitAuthentication()
    auth.InitCookieStore()
    database.InitDatabase()
    defer database.DeinitDatabase()

    // Setup various services
	router := mux.NewRouter()

	// Setup middleware
	router.Use(auth.AuthMiddleware)
    router.Use(mux.CORSMethodMiddleware(router))
    router.Use(auth.CORSMiddleware)

	log.Printf("[Router] Initialized router!")

    router.HandleFunc("/profile", handlers.HandleProfileQuery).Methods(http.MethodGet, http.MethodOptions)
    router.HandleFunc("/profile/{profile}", handlers.HandleProfile).Methods(http.MethodGet, http.MethodOptions)

    router.HandleFunc("/settings/profile", handlers.HandleProfileSettingsUpdate).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/auth/logout", handlers.HandleLogout).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}", handlers.HandleProviderLogin).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}/callback", handlers.HandleAuthCallback).Methods(http.MethodGet, http.MethodOptions)

	address := config.Envs.Host + ":" + config.Envs.Port

	log.Printf("[Router] Starting router...")

	// Serve static files
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	log.Printf("[Router] Listening on %s!", address)
	log.Fatalf("%s", http.ListenAndServe(address, router).Error())
}

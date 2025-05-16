package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/handlers"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/middleware"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Initializing piquel.fr API...\n")

	// Intialize services
	config.LoadConfig()
	auth.InitAuthentication()
	auth.InitCookieStore()
	database.InitDatabase()
	defer database.DeinitDatabase()
	types.Init()

	// Initialize the router
	router := mux.NewRouter()
	middleware.Setup(router)

	log.Printf("[Router] Initialized router!\n")

	router.HandleFunc("/profile", handlers.HandleBaseProfile).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/profile/{profile}", handlers.HandleProfile).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.HandleFunc("/auth/logout", handlers.HandleLogout).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}", handlers.HandleProviderLogin).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/{provider}/callback", handlers.HandleAuthCallback).Methods(http.MethodGet, http.MethodOptions)

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", http.ListenAndServe(address, router).Error())
}

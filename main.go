package main

import (
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/handlers"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/router"
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
    router := router.InitRouter()

    router.AddRoute("/profile", handlers.HandleProfileQuery, http.MethodGet)
    router.AddRoute("/profile/{profile}", handlers.HandleProfile, http.MethodGet)

    router.AddRoute("/settings/profile", handlers.HandleProfileSettingsUpdate, http.MethodPost)

	router.AddRoute("/auth/logout", handlers.HandleLogout, http.MethodGet)
	router.AddRoute("/auth/{provider}", handlers.HandleProviderLogin, http.MethodGet)
	router.AddRoute("/auth/{provider}/callback", handlers.HandleAuthCallback, http.MethodGet)

	address := config.Envs.Host + ":" + config.Envs.Port
    router.Start(address)
}

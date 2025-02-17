package main

import (
	"log"

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

    router.AddRoute("/profile", handlers.HandleProfileQuery, "GET")
    router.AddRoute("/profile/{profile}", handlers.HandleProfile, "GET")

    router.AddRoute("/settings/profile", handlers.HandleProfileSettingsUpdate, "POST")

	router.AddRoute("/auth/logout", handlers.HandleLogout, "GET")
	router.AddRoute("/auth/{provider}", handlers.HandleProviderLogin, "GET")
	router.AddRoute("/auth/{provider}/callback", handlers.HandleAuthCallback, "GET")

	address := config.Envs.Host + ":" + config.Envs.Port
    router.Start(address)
}

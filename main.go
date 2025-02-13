package main

import (
	"log"

	"github.com/PiquelChips/piquel.fr/handlers"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/router"
	"github.com/joho/godotenv"
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

	router.AddRoute("/", handlers.HandleHome, "GET")

    router.AddRoute("/profile", handlers.HandleBaseProfile, "GET")
    router.AddRoute("/profile/{profile}", handlers.HandleProfile, "GET")

    router.AddRoute("/settings", handlers.HandleSettingsRedirect, "GET")
    router.AddRoute("/settings/profile", handlers.HandleProfileSettings, "GET")
    router.AddRoute("/settings/profile", handlers.HandleProfileSettingsUpdate, "POST")

	// Basic public endpoints
	router.AddRoute("/minecraft", handlers.HandleMinecraft, "GET")
	router.AddRoute("/dirk", handlers.HandleDirk, "GET")

	// SNT endpoints
	router.AddRoute("/snt/linus", handlers.HandleLinus, "GET")
	router.AddRoute("/snt/linux", handlers.HandleLinux, "GET")
	router.AddRoute("/snt/musset", handlers.HandleMusset, "GET")

	// Auth
	router.AddRoute("/auth/login", handlers.HandleLogin, "GET")
	router.AddRoute("/auth/logout", handlers.HandleLogout, "GET")
	router.AddRoute("/auth/{provider}", handlers.HandleProviderLogin, "GET")
	router.AddRoute("/auth/{provider}/callback", handlers.HandleAuthCallback, "GET")

	godotenv.Load()
	address := config.Envs.Host + ":" + config.Envs.Port
    router.Start(address)
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/piquel-fr/api/api"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/email"
	"github.com/piquel-fr/api/services/users"
	gh "github.com/piquel-fr/api/utils/github"
	"github.com/piquel-fr/api/utils/oauth"
)

func main() {
	log.Printf("Initializing piquel.fr API...\n")

	// Intialize external services
	config.LoadConfig()
	gh.InitGithubClient()
	oauth.InitOAuth()
	database.InitDatabase()
	defer database.Connection.Close()

	userService := users.NewRealUserService()
	authService := auth.NewRealAuthService(userService)
	emailService := email.NewRealEmailService()

	config.UsernameBlacklist = userService.GetUsernameBlacklist()
	config.Policy = authService.GetPolicy()

	router, err := api.CreateRouter(userService, authService, emailService)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)
	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", server.ListenAndServe().Error())
}

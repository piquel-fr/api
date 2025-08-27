package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/handlers"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/docs"
	gh "github.com/piquel-fr/api/utils/github"
	"github.com/piquel-fr/api/utils/middleware"
)

func main() {
	log.Printf("Initializing piquel.fr API...\n")

	// Intialize external services
	config := config.LoadConfig()
	gh := gh.InitGithubClient(config)
	connection, queries := database.InitDatabase(config)
	defer connection.Close()

	handler := handlers.Handler{
		AuthService: auth.NewRealAuthService(config, queries),
		DocsService: docs.NewRealDocsService(gh),
	}

	address := fmt.Sprintf("0.0.0.0:%s", config.Envs.Port)

	server := http.Server{
		Addr: address,
		Handler: middleware.AddMiddleware(handler.CreateHttpHandler(config, queries, gh),
			middleware.CORSMiddleware,
		),
	}

	log.Printf("[Router] Starting router...\n")
	log.Printf("[Router] Listening on %s!\n", address)
	log.Fatalf("%s", server.ListenAndServe().Error())
}

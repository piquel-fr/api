package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Routes struct {
    Routes []Route
}

func loadRouteConfig() {
    // Load routes config
    routesData, err := os.ReadFile(fmt.Sprintf("%s/routes.yml", Envs.ConfigPath))
    if err != nil {
        panic(err)
    }
    err = yaml.Unmarshal(routesData, &Routes)
    if err != nil {
        panic(err)
    }

    // Fix up Routes data
    for _, route := range Routes.Routes {
        if route.Method == "" {
            route.Method = "GET"
        } else if route.Method != "GET" && route.Method != "POST" && route.Method != "DELETE" {
            log.Fatalf("Method %s of route %s is not valid!", route.Method, route.Name)
        }
    }

    log.Printf("[Config] Loaded routes config!")
}

package config

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

var routesConfig struct {
	routes     map[string]*Route
	slugRoutes []*Route
}

func loadRoutesConfig() {
	var configRoutes struct {
		Routes []*Route
	}

	routesConfig.routes = make(map[string]*Route)

	// Load routes config
	routesData, err := os.ReadFile(fmt.Sprintf("%s/routes.yml", Envs.ConfigPath))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(routesData, &configRoutes)
	if err != nil {
		panic(err)
	}

	// Fix up Routes data
	for _, route := range configRoutes.Routes {
		if route.Method == "" {
			route.Method = "GET"
		} else if route.Method != "GET" && route.Method != "POST" && route.Method != "DELETE" {
			log.Fatalf("Method %s of route %s is not valid!", route.Method, route.Name)
		}

		routesConfig.routes[route.Name] = route
		if route.Slug {
			routesConfig.slugRoutes = append(routesConfig.slugRoutes, route)
		}
	}

	log.Printf("[Config] Loaded routes config!")
}

func RouteExists(route string) bool {
	return routesConfig.routes[route] != nil
}

func GetRouteConfig(r *http.Request) *Route {
	route := routesConfig.routes[r.URL.Path]

	if route != nil {
		return route
	}

	vars := mux.Vars(r)
	for key := range vars {
		log.Printf("%s", key)
	}
	return route
}

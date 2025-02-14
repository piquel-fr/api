package router

import (
	"errors"
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/gorilla/mux"
)

type Router struct {
    Router *mux.Router
}

func InitRouter() *Router {
    router := &Router{Router: mux.NewRouter()}

    // Setup middleware
	router.Router.Use(auth.AuthMiddleware)

    log.Printf("[Router] Initialized router!")
    return router
}

func (router *Router) Start(address string) {
    log.Printf("[Router] Starting router...")

    // Serve static files
	router.Router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("public"))))

    log.Printf("[Router] Listening on %s!", address)
    log.Fatalf("%s", http.ListenAndServe(address, router.Router).Error())
}

func (router *Router) AddRoute(route string, handler func(http.ResponseWriter, *http.Request), method string) {
    if !config.RouteExists(route) {
        panic(errors.New("Please add %s route to config!"))
    }
    router.Router.HandleFunc(route, handler).Methods(method)
}

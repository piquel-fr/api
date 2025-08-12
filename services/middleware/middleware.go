package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(router *mux.Router) {
	router.Use(CORSMiddleware(router))
}

func CORSMiddleware(router *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				origin = r.Host
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "43100")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				/*
					var match mux.RouteMatch
					if !router.Match(r, &match) {
						http.NotFound(w, r)
						return
					}

					allMethods, err := match.Route.GetMethods()
					if err != nil {
						http.NotFound(w, r)
						panic(err)
					}

					w.Header().Set("Access-Control-Allow-Methods", strings.Join(allMethods, ","))
				*/
				w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

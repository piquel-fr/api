package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Setup(router *mux.Router) {
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(cORSMiddleware)
}

func cORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			origin = r.Host
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Max-Age", "43100")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			// Return immediately for OPTIONS requests
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

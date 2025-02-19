package middlewares

import (
	"net/http"
	"strings"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/gorilla/mux"
)

func Setup(router *mux.Router) {
	//router.Use(authMiddleware)
    router.Use(mux.CORSMethodMiddleware(router))
    router.Use(cORSMiddleware)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/auth/") {
            next.ServeHTTP(w, r)
            return
        }

        err := auth.VerifyUserSession(r)
        if err != nil {
            http.Error(w, "You are not authenticated", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
	})
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

package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PiquelChips/piquel.fr/services/config"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			if r.Host != r.Header.Get("Host") {
				http.Error(w, "Missing Origin Header", http.StatusForbidden)
				return
			}

			// Same origin
			next.ServeHTTP(w, r)
			return
		}

        isValidOrigin, allowCredentials := validateOrigin(origin, config.CORS.AllowedOrigins)

		if !isValidOrigin {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
		}

        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.CORS.MaxAge))
        if allowCredentials {
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }

		if r.Method == http.MethodOptions {
			// Return immediately for OPTIONS requests
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateOrigin(origin string, allowedOrigins map[string]bool) (bool, bool) {
	if origin == "" {
		return false, false
	}

	for allowed, credentials := range allowedOrigins {
		if allowed == origin {
			return true, credentials
		}

		// For example *.piquel.fr
		if strings.Contains(allowed, "*.") {
			// Would then be .piquel.fr
			domain := strings.Split(allowed, "*.")[1]
			if strings.HasSuffix(origin, domain) {
				return true, credentials
			}
		}
	}

	return false, false
}

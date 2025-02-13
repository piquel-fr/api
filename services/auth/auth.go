package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const SessionName = "user_session"

func InitAuthentication() {
	goth.UseProviders(
		google.New(
			config.Envs.GoogleClientID,
			config.Envs.GoogleClientSecret,
			buildCallbackURL("google"),
			"email", "profile",
		),
		github.New(
			config.Envs.GithubClientID,
			config.Envs.GithubClientSecret,
			buildCallbackURL("github"),
			"user:email",
		),
	)

	log.Printf("[Auth] Initialized auth service!\n")
}

func buildCallbackURL(provider string) string {
	var url string
	if config.Envs.SSL == "true" {
		url = fmt.Sprintf("https://%s/auth/%s/callback", config.Envs.PublicHost, provider)
	} else {
		url = fmt.Sprintf("http://%s:%s/auth/%s/callback", config.Envs.PublicHost, config.Envs.Port, provider)
	}
	log.Printf("[Auth] Added auth provider listener for %s on %s", provider, url)
	return url
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        /*
		routeConfig := config.RouteSettings[r.URL.Path]
		if !routeConfig.IsAuthenticated {
			next.ServeHTTP(w, r)
			return
		}

		session, err := GetSessionUser(r)
		if err != nil {
        http.Redirect(w, r, "/auth/login", http.StatusMethodNotAllowed)
			return
		}

        if routeConfig.Permission == "" {
		    next.ServeHTTP(w, r)
            return
        }

		user, err := database.Queries.GetUserByEmail(r.Context(), session.Email)
		if err != nil {
			if err != pgx.ErrNoRows {
				http.Redirect(w, r, "/auth/login", http.StatusMethodNotAllowed)
				return
			}
	        panic(err)
		}
        
        if permissions.UserHasPermission(user.ID, routeConfig.Permission) {
		    next.ServeHTTP(w, r)
            return
        }

        http.Redirect(w, r, "/auth/login", http.StatusMethodNotAllowed)
        */
        next.ServeHTTP(w, r)
	})
}

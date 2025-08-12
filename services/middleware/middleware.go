package middleware

import (
	"context"
	goErrors "errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
)

type Middleware func(http.Handler) http.Handler

// first in stack will run first
func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}

		return next
	}
}

func AddMiddleware(router http.Handler, middlewares ...Middleware) http.Handler {
	return CreateStack(middlewares...)(router)
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			origin = r.Host
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Max-Age", "43100")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.GetUserId(r)
		if err != nil {
			if !goErrors.Is(err, errors.ErrorNotAuthenticated) {
				errors.HandleError(w, r, err)
				return
			}
		}

		user, err := database.Queries.GetUserById(r.Context(), userId)
		if err != nil {
			if !goErrors.Is(err, pgx.ErrNoRows) {
				errors.HandleError(w, r, err)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", &user)))
	})
}

func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		if user == nil {
			http.Error(w, "please login to access this resource", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromRequest(r *http.Request) *repository.User {
	return r.Context().Value("user").(*repository.User)
}

func CreateOptionsHandler(methods ...string) http.Handler {
	methods = append(methods, "OPTIONS")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		w.WriteHeader(http.StatusOK)
	})
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/errors"
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}

func RequireAuthMiddleware(auth auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := auth.GetToken(r)
			if err != nil {
				errors.HandleError(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CreateOptionsHandler(methods ...string) http.Handler {
	methods = append(methods, "OPTIONS")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		w.WriteHeader(http.StatusOK)
	})
}

func AuthMiddleware(auth auth.AuthService, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.HandleHumaError(api, ctx, errors.ErrorNotAuthenticated)
			return
		}
		tokenString := parts[1]

		userId, err := auth.GetUserIdFromToken(ctx.Context(), tokenString)
		if err != nil {
			errors.HandleHumaError(api, ctx, err)
			return
		}

		next(huma.WithValue(ctx, "userId", userId))
	}
}

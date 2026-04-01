package middleware

import (
	"context"
	"net/http"
	"strings"

	"smartqueue/internal/service"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			claims, err := auth.ParseToken(tokenString)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

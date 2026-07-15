package middleware

import (
	"context"
	"net/http"
	"strings"
	"task-manager-go/internal/ctxkey"
	"task-manager-go/services"

	"github.com/golang-jwt/jwt/v5"
)

// Auth is a middleware factory that returns a middleware for validating JWT tokens.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	secretBytes := []byte(jwtSecret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims := &services.Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secretBytes, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxkey.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

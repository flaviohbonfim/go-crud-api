package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-crud-api/internal/config"
	customhttp "go-crud-api/pkg/web"
	"go-crud-api/pkg/jwt"

	"github.com/rs/zerolog/log"
)

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const ( // Define context keys
	ContextKeyUserID contextKey = "userID"
	ContextKeyRole   contextKey = "role"
)

// AuthMiddleware validates JWT tokens and adds user info to context.
func AuthMiddleware(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				customhttp.RespondWithError(w, "unauthorized", "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				customhttp.RespondWithError(w, "unauthorized", "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := jwt.ValidateToken(tokenString, cfg.JWTSecret)
			if err != nil {
				log.Error().Err(err).Msg("Invalid JWT token")
				customhttp.RespondWithError(w, "unauthorized", "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user ID and role to context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// HasRoleMiddleware checks if the authenticated user has one of the required roles.
func HasRoleMiddleware(requiredRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok {
				customhttp.RespondWithError(w, "forbidden", "Role not found in context", http.StatusForbidden)
				return
			}

			roleFound := false
			for _, role := range requiredRoles {
				if userRole == role {
					roleFound = true
					break
				}
			}

			if !roleFound {
				customhttp.RespondWithError(w, "forbidden", "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

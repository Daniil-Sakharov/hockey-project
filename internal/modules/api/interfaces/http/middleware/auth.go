package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// ContextKey is a custom type for context keys.
type ContextKey string

const (
	// UserContextKey is the key for user claims in context.
	UserContextKey ContextKey = "user"
)

// AuthMiddleware provides JWT authentication middleware.
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAuth is middleware that requires a valid JWT token.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"unauthorized","message":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":"unauthorized","message":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.authService.ValidateAccessToken(tokenString)
		if err != nil {
			logger.Debug(ctx, "Invalid token", zap.Error(err))
			http.Error(w, `{"error":"unauthorized","message":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx = context.WithValue(ctx, UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth is middleware that optionally validates JWT token.
// If token is present and valid, user claims are added to context.
// If token is missing or invalid, request proceeds without user claims.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString := parts[1]
				claims, err := m.authService.ValidateAccessToken(tokenString)
				if err == nil {
					ctx = context.WithValue(ctx, UserContextKey, claims)
				}
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireSubscription is middleware that requires a specific subscription tier.
func (m *AuthMiddleware) RequireSubscription(tier string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetUserFromContext(r.Context())
			if claims == nil {
				http.Error(w, `{"error":"unauthorized","message":"authentication required"}`, http.StatusUnauthorized)
				return
			}

			// Check subscription tier
			if !hasRequiredTier(claims.SubscriptionTier, tier) {
				http.Error(w, `{"error":"forbidden","message":"upgrade your subscription to access this feature"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext retrieves user claims from context.
func GetUserFromContext(ctx context.Context) *services.JWTClaims {
	claims, ok := ctx.Value(UserContextKey).(*services.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}

// hasRequiredTier checks if user's tier meets the required tier.
func hasRequiredTier(userTier, requiredTier string) bool {
	tierLevels := map[string]int{
		"free":  0,
		"pro":   1,
		"ultra": 2,
	}

	userLevel, ok := tierLevels[userTier]
	if !ok {
		return false
	}

	requiredLevel, ok := tierLevels[requiredTier]
	if !ok {
		return false
	}

	return userLevel >= requiredLevel
}

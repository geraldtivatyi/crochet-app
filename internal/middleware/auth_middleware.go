package middleware

import (
	"net/http"
	"strings"

	"Crochet/internal/auth"
)

// RequireAuth checks for a valid Bearer token in the Authorization header
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Parse "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Simulate token validation (In production, parse JWT or check database/Redis)
		user, err := validateToken(token)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Attach user claims to request context!
		ctx := auth.WithUser(r.Context(), user)

		// Pass the request WITH the new context to the next handler in chain
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (auth.UserClaims, error) {
	// Dummy token logic for illustration
	if token == "secret-artisan-token" {
		return auth.UserClaims{
			UserID:  "user_123",
			Role:    "artisan",
			StoreID: "store_cape_town",
		}, nil
	}
	return auth.UserClaims{}, http.ErrNoCookie
}

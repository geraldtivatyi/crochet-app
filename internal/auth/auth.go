package auth

import (
	"context"
	"errors"
)

// UserClaims holds authenticated user info
type UserClaims struct {
	UserID  string
	Role    string // e.g., "admin", "artisan"
	StoreID string
}

// Private type for context key to avoid collisions
type contextKey string

const UserContextKey contextKey = "userClaims"

// Helper: Attach UserClaims to context
func WithUser(ctx context.Context, user UserClaims) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// Helper: Extract UserClaims from context in handlers
func UserFromContext(ctx context.Context) (UserClaims, error) {
	user, ok := ctx.Value(UserContextKey).(UserClaims)
	if !ok {
		return UserClaims{}, errors.New("no user found in context")
	}
	return user, nil
}

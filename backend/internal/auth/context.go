package auth

import "context"

type UserContext struct {
	UserID   string
	Username string
	Role     string
	StoreID  *string
}

type contextKey string

const UserContextKey contextKey = "user"

func WithUser(ctx context.Context, user *UserContext) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

func GetUser(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
}

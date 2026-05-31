package auth

import "context"

type contextKey string

const UserContextKey contextKey = "user"

func GetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(*User)
	return user, ok
}

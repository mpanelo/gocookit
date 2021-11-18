package context

import (
	"context"

	"github.com/mpanelo/gocookit/models"
)

type contextKey string

const (
	userKey = contextKey("user")
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if value := ctx.Value(userKey); value != nil {
		user, ok := value.(*models.User)
		if ok {
			return user
		}
		return nil
	}
	return nil
}

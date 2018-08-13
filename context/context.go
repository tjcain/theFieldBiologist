package context

import (
	"context"

	"github.com/tjcain/theFieldBiologist/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// WithUser ...
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User ...
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}

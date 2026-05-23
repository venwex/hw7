package auth

import (
	"context"
	"errors"
)

type contextKey string

const currentUserKey contextKey = "current_user"

var ErrUnauthorized = errors.New("unauthorized")

func WithCurrentUser(ctx context.Context, user CurrentUser) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

func GetCurrentUser(ctx context.Context) (CurrentUser, error) {
	user, ok := ctx.Value(currentUserKey).(CurrentUser)
	if !ok {
		return CurrentUser{}, ErrUnauthorized
	}

	return user, nil
}

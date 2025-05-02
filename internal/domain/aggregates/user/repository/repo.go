package repository

import (
	"context"
	vo "physk/internal/domain/aggregates/user/value_objects"
)

type UserRepository interface {
	UsernameExists(ctx context.Context, username vo.UserName) (bool, error)
	LoginExists(ctx context.Context, login vo.UserLogin) (bool, error)
	EmailExists(ctx context.Context, email vo.UserEmail) (bool, error)
}

/// todo: add provider

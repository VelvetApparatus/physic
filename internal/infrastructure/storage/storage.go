package storage

import (
	"context"
	"errors"
	"physk/internal/domain/aggregates/user"
)

var (
	ErrorNotFound = errors.New("not found")
)

type Storage interface {
	CreateUser(ctx context.Context, user user.User) error
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
}

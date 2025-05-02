package usecase

import (
	"context"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/storage"
)

type WriteModel interface {
	Create(ctx context.Context, u user.User) error
}

type ReadModel interface {
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
}

type UseCase struct {
	s storage.Storage
}

func (u *UseCase) Create(ctx context.Context, us user.User) error {
	return u.s.CreateUser(ctx, us)
}

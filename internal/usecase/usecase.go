package usecase

import (
	"context"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/storage"
)

type WriteModel interface {
	Create(ctx context.Context, u user.User) error
}

type ReadModel interface {
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (user.User, error)
}

type UseCase struct {
	s storage.Storage
}

func NewUseCase(st storage.Storage) *UseCase {
	return &UseCase{s: st}
}

func (u *UseCase) Create(ctx context.Context, us user.User) error {
	return u.s.CreateUser(ctx, us)
}

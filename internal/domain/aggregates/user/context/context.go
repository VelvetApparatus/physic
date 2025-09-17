package context

import (
	"context"
	"physk/internal/domain/aggregates/user/repository"
)

var (
	repo repository.UserRepository
)

type UserCtx struct {
	context.Context

	Repo repository.UserRepository
}

func InitUserContextFabric(
	r repository.UserRepository,
) {
	repo = r
}

// todo: or make it like context.Context pattern?

func NewUserContext(ctx context.Context) UserCtx {
	return UserCtx{
		Repo:    repo,
		Context: ctx,
	}
}

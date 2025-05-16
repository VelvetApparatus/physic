package user

import (
	"go.uber.org/fx"
	userCtx "physk/internal/domain/aggregates/user/context"
	userRepo "physk/internal/domain/aggregates/user/repository"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"user-aggregate",
		userRepo.ProvideModule(),
		userCtx.ProvideModule(),
	)
}

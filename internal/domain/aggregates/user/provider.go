package user

import (
	"go.uber.org/fx"
	userCtx "physk/internal/domain/aggregates/user/context"
	postgresRepo "physk/internal/domain/aggregates/user/repository/postgres"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"user-aggregate",
		postgresRepo.ProvideModule(),
		userCtx.ProvideModule(),
	)
}

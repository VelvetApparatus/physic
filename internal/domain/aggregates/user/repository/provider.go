package repository

import (
	"go.uber.org/fx"
	"physk/internal/domain/aggregates/user/repository/postgres"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"user-repository",
		fx.Provide(postgres.NewUserRepository),
	)
}

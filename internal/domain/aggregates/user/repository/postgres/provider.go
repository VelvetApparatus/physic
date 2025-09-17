package postgres

import (
	"go.uber.org/fx"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"user-repository",
		fx.Provide(NewUserRepository),
	)
}

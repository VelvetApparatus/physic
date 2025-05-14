package infrastructure

import (
	"go.uber.org/fx"
	"physk/internal/infrastructure/db"
	"physk/internal/infrastructure/storage/postgres"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"infrastructure",
		fx.Provide(
			db.NewConnection,
			postgres.NewStorage,
		),
	)
}

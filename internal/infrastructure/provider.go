package infrastructure

import (
	"go.uber.org/fx"
	"physk/internal/infrastructure/db"
	"physk/internal/infrastructure/db/storage/postgres"
	"physk/internal/infrastructure/delivery/fiber"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"infrastructure",

		fx.Provide(
			db.NewConnection,
			postgres.NewStorage,
		),

		fiber.ProvideModule(),
	)
}

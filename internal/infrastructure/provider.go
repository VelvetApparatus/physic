package infrastructure

import (
	"go.uber.org/fx"
	"physk/internal/infrastructure/db"
	"physk/internal/infrastructure/db/storage/postgres"
	"physk/internal/infrastructure/delivery/fiber"
	"physk/internal/infrastructure/s3"
	"physk/internal/infrastructure/s3/storage/minio"
)

func ProvideModule() fx.Option {
	return fx.Module(
		"infrastructure",

		fx.Provide(
			db.NewConnection,
			s3.NewMinioClient,
			postgres.NewStorage,
			minio.NewS3Storage,
		),

		fiber.ProvideModule(),
	)
}

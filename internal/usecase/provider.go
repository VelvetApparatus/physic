package usecase

import "go.uber.org/fx"

func ProvideModule() fx.Option {
	return fx.Module(
		"use-cases",
		fx.Provide(
			NewUseCase,
			fx.Annotate(
				NewUseCase,
				fx.As(new(WriteModel)),
				fx.As(new(ReadModel))),
		),
	)
}

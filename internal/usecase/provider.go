package usecase

import "go.uber.org/fx"

func ProvideModule() fx.Option {
	return fx.Module(
		"use-cases",
		fx.Provide(NewUseCase),
	)
}

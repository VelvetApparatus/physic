package usecase

import "go.uber.org/fx"

func ProvideUseCases() fx.Option {
	return fx.Module("use-cases")
}

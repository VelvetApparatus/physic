package controller

import "go.uber.org/fx"

func ProvideModule() fx.Option {
	return fx.Module(
		"controller",
		fx.Provide(NewController),
	)
}

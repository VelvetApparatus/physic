package controller

import "go.uber.org/fx"

func ProvideController() fx.Option {
	return fx.Module(
		"controller",
		fx.Invoke(NewController),
	)
}

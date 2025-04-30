package context

import "go.uber.org/fx"

func ProvideUserContext() fx.Option {
	return fx.Module(
		"user-context",
		fx.Provide(
			InitUserContextFabric,
		),
	)
}

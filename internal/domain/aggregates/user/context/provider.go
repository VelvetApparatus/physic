package context

import "go.uber.org/fx"

func ProvideUserContext() fx.Option {
	return fx.Module(
		"user-context-fabric",
		fx.Invoke(
			InitUserContextFabric,
		),
	)
}

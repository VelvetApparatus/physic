package context

import "go.uber.org/fx"

func ProvideModule() fx.Option {
	return fx.Module(
		"user-context-fabric",
		fx.Invoke(
			InitUserContextFabric,
		),
	)
}

package fiber

import "go.uber.org/fx"

func ProvideRouter() fx.Option {
	return fx.Module(
		"router",
		fx.Provide(NewRouter),
		fx.Invoke(StartFiberRouter),
	)

}

package fiber

import "go.uber.org/fx"

func ProvideModule() fx.Option {
	return fx.Module(
		"router",
		fx.Provide(NewRouter),
		fx.Invoke(StartFiberRouter),
	)

}

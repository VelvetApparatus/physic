package app

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"os"
	"physk/internal/config"
	"physk/internal/controller"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure"
	"physk/internal/usecase"
)

func App(ctx context.Context) {
	app := fx.New(
		fx.Module("check",
			fx.Provide(func() context.Context { return ctx }),
			fx.Invoke(config.Init),
		),

		infrastructure.ProvideModule(),

		user.ProvideModule(),

		usecase.ProvideModule(),

		controller.ProvideModule(),
	)

	if err := app.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err.Error())
		return
	}

	app.Run()

}

package app

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"os"
	"physk/internal/config"
	"physk/internal/controller"
	userCtx "physk/internal/domain/aggregates/user/context"
	"physk/internal/infrastructure"
	"physk/internal/infrastructure/delivery/fiber"
	"physk/internal/usecase"
)

func App(ctx context.Context) {
	app := fx.New(
		fx.Supply(ctx),
		fx.Invoke(config.Init),

		infrastructure.ProvideModule(),

		userCtx.ProvideModule(),

		usecase.ProvideModule(),

		controller.ProvideModule(),

		fiber.ProvideModule(),
	)

	if err := app.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err.Error())
		return
	}

	app.Run()

}

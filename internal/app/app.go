package app

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"os"
	"physk/internal/controller"
	userCtx "physk/internal/domain/aggregates/user/context"
	"physk/internal/infrastructure/delivery/fiber"
	"physk/internal/usecase"
)

func App(ctx context.Context) {
	app := fx.New(
		fx.Supply(ctx),
		//fx.config ?

		userCtx.ProvideUserContext(),

		usecase.ProvideUseCases(),

		controller.ProvideController(),

		fiber.ProvideRouter(),
	)

	if err := app.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err.Error())
		return
	}

	app.Run()

}

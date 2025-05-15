package app

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"os"
	"physk/internal/config"
	"physk/internal/controller"
	userCtx "physk/internal/domain/aggregates/user/context"
	userRepo "physk/internal/domain/aggregates/user/repository"
	"physk/internal/infrastructure"
	"physk/internal/usecase"
)

func App(ctx context.Context) {
	app := fx.New(
		fx.Supply(ctx),
		fx.Invoke(config.Init),

		infrastructure.ProvideModule(),

		// merge to userAggregate
		userCtx.ProvideModule(),

		userRepo.ProvideModule(),

		usecase.ProvideModule(),

		controller.ProvideModule(),
	)

	if err := app.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err.Error())
		return
	}

	app.Run()

}

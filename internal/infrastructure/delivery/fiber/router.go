package fiber

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"physk/internal/config"
	"physk/internal/infrastructure/delivery/fiber/mdws"
	tokenServ "physk/internal/services/token_service"

	"physk/internal/controller"
	"time"
)

type Router struct {
	ctrl *controller.Controller
}

func NewRouter(ctrl *controller.Controller) *Router {
	return &Router{ctrl: ctrl}
}

func (r *Router) MapRoutes(
	group fiber.Router,
) {
	// middlewares
	authedMiddleware := mdws.TokenValidationMDW(tokenServ.NewTokenService())

	// commands
	group.Post("/register", r.Register())
	group.Post("/login", r.Login())

	// queries
	group.Get("/me", authedMiddleware, r.GetMe())

}

func StartFiberRouter(router *Router, lc fx.Lifecycle) error {
	app := fiber.New(
		&fiber.Settings{
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 15,
			IdleTimeout:  time.Minute,
		},
	)

	recoverConfig := fiberRecover.ConfigDefault
	recoverConfig.EnableStackTrace = true
	app.Use(fiberRecover.New(recoverConfig))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	v1Group := app.Group("api/v1")
	router.MapRoutes(v1Group)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := app.Listen(fmt.Sprintf("%s:%d", config.C().App.Host, config.C().App.Port))
			if err != nil {
				return fmt.Errorf("listen: %w", err)
			}
			return nil

		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
	return nil
}

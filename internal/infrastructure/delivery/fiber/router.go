package fiber

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
	isAdminMiddleware := mdws.IsAdminRole()

	// subgroups
	userGroup := group.Group("/user")
	collectionGroup := group.Group("/collection")

	// commands
	userGroup.Post("/register", r.Register())
	userGroup.Post("/login", r.Login())
	collectionGroup.Post("/create", authedMiddleware, isAdminMiddleware, r.CreateCollection())
	collectionGroup.Post("/attach", authedMiddleware, isAdminMiddleware, r.AttachImageToCollection())

	// queries
	userGroup.Get("/me", authedMiddleware, r.GetMe())
	collectionGroup.Post("/ids", r.GetImageIDsByCollectionID())
	collectionGroup.Post("/image", r.GetImageByID())

}

func StartFiberRouter(lc fx.Lifecycle, router *Router) error {
	app := fiber.New(
		fiber.Config{
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
			addr := fmt.Sprintf("%s:%d", config.C().App.Host, config.C().App.Port)
			fmt.Printf("Starting Fiber app on %s\n", addr)
			go app.Listen(addr)
			//if err != nil {
			//	return fmt.Errorf("failed to start Fiber app: %w", err)
			//}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down Fiber app")
			return app.Shutdown()
		},
	})
	return nil
}

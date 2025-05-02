package fiber

import (
	"github.com/gofiber/fiber"
	"physk/internal/controller"
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

	// commands
	group.Post("/register", r.Register())
	group.Post("/login", r.Login())

	// queries
	group.Get("/me", r.GetMe())

}

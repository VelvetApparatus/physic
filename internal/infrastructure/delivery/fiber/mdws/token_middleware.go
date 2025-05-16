package mdws

import (
	"errors"
	"github.com/gofiber/fiber"
	tokenServ "physk/internal/services/token_service"
	"strings"
)

func TokenValidationMDW(service tokenServ.TokenService) fiber.Handler {

	return func(fiberCtx *fiber.Ctx) {
		parts := strings.Split(fiberCtx.Get("Authorization"), ": ")
		if len(parts) != 1 {
			fiberCtx.Status(fiber.StatusUnauthorized)
		}

		token := parts[0]

		claims, err := service.ValidateToken(token)
		if err != nil {
			if errors.Is(err, tokenServ.InvalidTokenError) {
				fiberCtx.SendStatus(fiber.StatusUnauthorized)
			}
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
		}
		fiberCtx.Locals("claims", claims)
		fiberCtx.Next()
	}

}

func IsAdminRole() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		claims, ok := fiberCtx.Locals("claims").(*tokenServ.Claims)
		if !ok {
			fiberCtx.SendStatus(fiber.StatusUnauthorized)
			return
		}

		if !claims.Role.IsAdmin() {
			fiberCtx.SendStatus(fiber.StatusUnauthorized)
			return
		}

		fiberCtx.Next()
	}
}

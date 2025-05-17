package mdws

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	tokenServ "physk/internal/services/token_service"
	"strings"
)

func TokenValidationMDW(service tokenServ.TokenService) fiber.Handler {

	return func(fiberCtx *fiber.Ctx) error {
		parts := strings.Split(fiberCtx.Get("Authorization"), ": ")
		if len(parts) != 1 {
			fiberCtx.Status(fiber.StatusUnauthorized)
		}

		token := parts[0]

		claims, err := service.ValidateToken(token)
		if err != nil {
			if errors.Is(err, tokenServ.InvalidTokenError) {
				return fiberCtx.SendStatus(fiber.StatusUnauthorized)
			}
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)
		}
		fiberCtx.Locals("claims", claims)
		return fiberCtx.Next()
	}

}

func IsAdminRole() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		claims, ok := fiberCtx.Locals("claims").(*tokenServ.Claims)
		if !ok {
			return fiberCtx.SendStatus(fiber.StatusUnauthorized)
		}

		if !claims.Role.IsAdmin() {
			return fiberCtx.SendStatus(fiber.StatusUnauthorized)

		}

		return fiberCtx.Next()
	}
}

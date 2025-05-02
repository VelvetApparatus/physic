package fiber

import (
	"context"
	"errors"
	"github.com/gofiber/fiber"
	"log/slog"
	"physk/internal/controller"
	userErrors "physk/internal/domain/aggregates/user/errors"
	"physk/internal/infrastructure/delivery/dto"
	"physk/internal/infrastructure/delivery/fiber/models"
)

func (r *Router) Register() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.UserCreateRequest{}

		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("body", fiberCtx.Body()),
				slog.String("error", err.Error()),
			)
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		user, err := r.ctrl.RegisterUser(ctx, request)
		if err != nil {
			var userErr userErrors.UserError
			if errors.As(err, &userErr) {
				fiberCtx.Status(fiber.StatusBadRequest).JSON(models.InvalidRequestResponse{Msg: userErr.Error()})
				return
			}
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.UserCreateResponse{
			ID:       user.ID,
			Role:     user.Role.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
		}

		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "json response",
				slog.String("error", err.Error()),
				slog.Any("response", response),
			)
		}

	}
}

func (r *Router) Login() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.UserLoginRequest{}

		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parse",
				slog.String("error", err.Error()),
				slog.String("body", fiberCtx.Body()),
			)

			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		token, err := r.ctrl.Login(ctx, request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "login",
				slog.String("error", err.Error()),
			)

			if errors.Is(err, controller.UserNotFound) {
				fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid credentials")
				return
			}
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.UserLoginResponse{Token: token}
		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "send login response",
				slog.Any("response", response),
				slog.String("error", err.Error()),
			)
		}
	}
}

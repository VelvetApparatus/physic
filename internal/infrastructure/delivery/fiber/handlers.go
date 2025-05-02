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
	tokenServ "physk/internal/services/token_service"
	"strings"
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

func (r *Router) GetMe() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		parts := strings.Split(fiberCtx.Get("Authorization"), ": ")
		if len(parts) != 2 {
			fiberCtx.Status(fiber.StatusUnauthorized)
		}

		token := parts[1]

		user, err := r.ctrl.GetMe(ctx, token)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "get-me",
				slog.String("error", err.Error()),
			)

			// todo: do not catch errors from token service
			if errors.Is(err, tokenServ.InvalidTokenError) {
				fiberCtx.Status(fiber.StatusUnauthorized).SendString("invalid token")
				return
			}
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.GetMeResponse{
			ID:       user.ID,
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Role:     user.Role.String(),
		}

		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		slog.LogAttrs(
			ctx, slog.LevelError, "json response",
			slog.String("error", err.Error()),
			slog.Any("response", response),
		)

	}
}

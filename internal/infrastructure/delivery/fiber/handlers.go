package fiber

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"mime/multipart"
	"net/textproto"
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
			return
		}

		claims, ok := fiberCtx.Locals("claims").(tokenServ.Claims)
		if !ok {
			fiberCtx.SendStatus(fiber.StatusBadRequest)
			return
		}

		user, err := r.ctrl.GetMe(ctx, claims.UserID)
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

func (r *Router) CreateCollection() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := new(dto.CreateCollectionRequest)

		err := fiberCtx.BodyParser(request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parse",
				slog.String("handler", "create-collection"),
				slog.String("error", err.Error()),
			)

			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		collection, err := r.ctrl.CreateCollection(ctx, *request)
		if err != nil {
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.CreateCollectionResponse{ID: collection.ID}

		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "fiber: json",
				slog.String("handler", "create-collection"),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (r *Router) AttachImageToCollection() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		form, err := fiberCtx.MultipartForm()
		if err != nil {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		rawID, ok := form.Value["id"]
		if !ok || len(rawID) != 1 {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("no collection id")
			return
		}

		collectionID, err := uuid.Parse(rawID[0])
		if err != nil {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid collection id")
			return
		}

		formNames, ok := form.Value["name"]
		if !ok || len(formNames) != 1 {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("no image name")
			return
		}

		name := formNames[0]
		if len(name) == 0 {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid image name")
			return
		}

		formFiles, ok := form.File["image"]
		if !ok || len(formFiles) != 1 {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("no image file")
			return
		}

		fileHeader := formFiles[0]

		contentType := fileHeader.Header.Get("Content-Type")
		if len(contentType) == 0 {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid content-type")
			return
		}

		if contentType != "image/png" && contentType != "image/jpg" && contentType != "image/jpeg" {
			fiberCtx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid image content-type value: %s", contentType))
			return
		}

		f, err := fileHeader.Open()
		if err != nil {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid file")
			return
		}
		defer f.Close()

		bytes, err := io.ReadAll(f)
		if err != nil {
			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid file")
			return
		}

		request := dto.AddImageToCollectionRequest{
			CollectionID: collectionID,
			Name:         name,
			ContentType:  contentType,
			Data:         bytes,
		}

		imgID, err := r.ctrl.AddImageToCollection(ctx, request)

		if err != nil {
			if errors.Is(err, controller.CollectionNotFound) {
				fiberCtx.Status(fiber.StatusNotFound).SendString("not found")
				return
			}

			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.AddImageToCollectionResponse{ID: imgID}

		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "fiber: json",
				slog.String("handler", "add-image-to-collection"),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (r *Router) GetImageIDsByCollectionID() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.GetImageIdsByCollectionIDRequest{}
		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("handler", "get-images-ids-by-collection-id"),
				slog.String("error", err.Error()),
			)

			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		ids, err := r.ctrl.GetImageIDsByCollectionID(ctx, request.CollectionID)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "get-image-ids-by-collection-id",
				slog.String("error", err.Error()),
				slog.String("collectionID", request.CollectionID.String()),
			)

			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		response := dto.GetImagesIDsByCollectionIDResponse{ImageIDs: ids}
		err = fiberCtx.Status(fiber.StatusOK).JSON(response)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "fiber: json",
				slog.String("handler", "get-images-ids-by-collection-id"),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (r *Router) GetImageByID() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.GetImageByID{}
		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("handler", "get-image-by-id"),
				slog.String("error", err.Error()),
			)

			fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
			return
		}

		img, err := r.ctrl.GetImageByID(ctx, request.ID)
		if err != nil {
			if errors.Is(err, controller.ImageNotFound) {
				fiberCtx.SendStatus(fiber.StatusNotFound)
				return
			}
		}

		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		defer writer.Close()

		h := make(textproto.MIMEHeader)

		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				img.Name, img.Name))

		h.Set("Content-Type", img.ContentType)

		part, err := writer.CreatePart(h)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "create part",
				slog.String("error", err.Error()),
			)
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}
		_, err = part.Write(img.ImageData)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "write file",
				slog.String("error", err.Error()),
			)
			fiberCtx.SendStatus(fiber.StatusInternalServerError)
			return
		}

		fiberCtx.Set(fiber.HeaderContentType, writer.FormDataContentType())

		fiberCtx.SendBytes(buf.Bytes())
	}
}

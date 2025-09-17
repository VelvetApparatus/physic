package fiber

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"

	"github.com/google/uuid"
	"io"
	"log/slog"
	"physk/internal/controller"
	userErrors "physk/internal/domain/aggregates/user/errors"
	"physk/internal/infrastructure/delivery/dto"
	"physk/internal/infrastructure/delivery/fiber/models"
	tokenServ "physk/internal/services/token_service"
)

func (r *Router) Register() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.UserCreateRequest{}

		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("body", string(fiberCtx.Body())),
				slog.String("error", err.Error()),
			)
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")
		}

		user, err := r.ctrl.RegisterUser(ctx, request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "register user uc",
				slog.String("error", err.Error()),
			)
			var userErr userErrors.UserError
			if errors.As(err, &userErr) {
				return fiberCtx.Status(fiber.StatusBadRequest).JSON(models.InvalidRequestResponse{Msg: userErr.Error()})

			}
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.UserCreateResponse{
			ID:       user.ID,
			Role:     user.Role.String(),
			Username: user.Username.String(),
			Email:    user.Email.String(),
		}

		return fiberCtx.Status(fiber.StatusOK).JSON(response)
	}
}

func (r *Router) Login() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.UserLoginRequest{}

		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parse",
				slog.String("error", err.Error()),
				slog.String("body", string(fiberCtx.Body())),
			)

			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")

		}

		token, err := r.ctrl.Login(ctx, request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "login",
				slog.String("error", err.Error()),
			)

			if errors.Is(err, controller.UserNotFound) {
				return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid credentials")

			}
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.UserLoginResponse{Token: token}
		return fiberCtx.Status(fiber.StatusOK).JSON(response)
	}
}

func (r *Router) GetMe() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		claims, ok := fiberCtx.Locals("claims").(*tokenServ.Claims)
		if !ok {
			return fiberCtx.SendStatus(fiber.StatusUnauthorized)

		}

		user, err := r.ctrl.GetMe(ctx, claims.UserID)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "get-me",
				slog.String("error", err.Error()),
			)

			// todo: do not catch errors from token service
			if errors.Is(err, tokenServ.InvalidTokenError) {
				return fiberCtx.Status(fiber.StatusUnauthorized).SendString("invalid token")

			}
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.GetMeResponse{
			ID:       user.ID,
			Username: user.Username.String(),
			Email:    user.Email.String(),
			Role:     user.Role.String(),
		}

		return fiberCtx.Status(fiber.StatusOK).JSON(response)
	}
}

func (r *Router) CreateCollection() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
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

			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")

		}

		collection, err := r.ctrl.CreateCollection(ctx, *request)
		if err != nil {
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.CreateCollectionResponse{ID: collection.ID}

		return fiberCtx.Status(fiber.StatusOK).JSON(response)
	}
}

func (r *Router) AttachImageToCollection() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		form, err := fiberCtx.MultipartForm()
		if err != nil {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")

		}

		rawID, ok := form.Value["id"]
		if !ok || len(rawID) != 1 {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("no collection id")

		}

		collectionID, err := uuid.Parse(rawID[0])
		if err != nil {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid collection id")

		}

		formNames, ok := form.Value["name"]
		if !ok || len(formNames) != 1 {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("no image name")

		}

		name := formNames[0]
		if len(name) == 0 {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid image name")
		}

		var isPreview bool
		isPreviewStr, ok := form.Value["is_preview"]
		if !ok || len(isPreviewStr) != 1 {
			isPreview = false
		} else {
			isPreview, err = strconv.ParseBool(isPreviewStr[0])
			if err != nil {
				isPreview = false
			}
		}

		formFiles, ok := form.File["image"]
		if !ok || len(formFiles) != 1 {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("no image file")
		}

		fileHeader := formFiles[0]

		contentType := fileHeader.Header.Get("Content-Type")
		if len(contentType) == 0 {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid content-type")
		}

		if contentType != "image/png" && contentType != "image/jpg" && contentType != "image/jpeg" {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid image content-type value: %s", contentType))
		}

		f, err := fileHeader.Open()
		if err != nil {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid file")

		}
		defer f.Close()

		bytes, err := io.ReadAll(f)
		if err != nil {
			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid file")

		}

		request := dto.AddImageToCollectionRequest{
			CollectionID: collectionID,
			Name:         name,
			ContentType:  contentType,
			Data:         bytes,
			IsPreview:    isPreview,
		}

		imgID, err := r.ctrl.AddImageToCollection(ctx, request)

		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "add image to collection",
				slog.String("error", err.Error()),
			)
			if errors.Is(err, controller.CollectionNotFound) {
				return fiberCtx.Status(fiber.StatusNotFound).SendString("not found")

			}

			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.AddImageToCollectionResponse{ID: imgID}

		return fiberCtx.Status(fiber.StatusOK).JSON(response)
	}
}

func (r *Router) GetImageIDsByCollectionID() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
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

			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")

		}

		slog.LogAttrs(
			ctx, slog.LevelInfo, "request",
			slog.String("handler", "get image ids for collection"),
			slog.String("collection_id", request.CollectionID.String()),
		)

		ids, err := r.ctrl.GetImageIDsByCollectionID(ctx, request.CollectionID)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "get-image-ids-by-collection-id",
				slog.String("error", err.Error()),
				slog.String("collectionID", request.CollectionID.String()),
			)

			return fiberCtx.SendStatus(fiber.StatusInternalServerError)

		}

		response := dto.GetImagesIDsByCollectionIDResponse{ImageIDs: ids}
		return fiberCtx.Status(fiber.StatusOK).JSON(response)

	}
}

func (r *Router) GetImageByID() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.GetImageByID{}
		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("handler", "get-image-by-id"),
				slog.String("body", string(fiberCtx.Body())),
				slog.String("error", err.Error()),
			)

			return fiberCtx.Status(fiber.StatusBadRequest).SendString("invalid body")

		}

		img, err := r.ctrl.GetImageByID(ctx, request.ID)
		if err != nil {
			if errors.Is(err, controller.ImageNotFound) {
				return fiberCtx.SendStatus(fiber.StatusNotFound)

			}
		}

		fiberCtx.Set(fiber.HeaderContentType, img.ContentType)
		return fiberCtx.Send(img.ImageData)
	}
}

func (r *Router) GetCollections() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		collections, err := r.ctrl.GetCollections(ctx)
		if err != nil {
			return fiberCtx.SendStatus(fiber.StatusInternalServerError)
		}

		response := dto.GetCollections{}

		for _, colItem := range collections.Items {
			response.Items = append(response.Items, dto.GetCollectionItem{
				CollectionID:   colItem.CollectionID,
				PreviewImageID: colItem.PreviewID,
				CollectionName: colItem.Name,
			})
		}

		return fiberCtx.JSON(response)
	}
}

func (r *Router) DeleteCollection() fiber.Handler {
	return func(fiberCtx *fiber.Ctx) error {
		ctx, cancel := context.WithCancel(fiberCtx.Context())
		defer cancel()

		request := dto.DeleteCollectionRequest{}

		err := fiberCtx.BodyParser(&request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "body-parser",
				slog.String("handler", "delete_collection"),
				slog.String("error", err.Error()),
				slog.String("body", string(fiberCtx.Body())),
			)

			return fiberCtx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = r.ctrl.DeleteCollection(ctx, request)
		if err != nil {
			slog.LogAttrs(
				ctx, slog.LevelError, "delete collection",
				slog.String("error", err.Error()),
				slog.String("collectionID", request.CollectionID.String()),
			)

			return fiberCtx.SendStatus(fiber.StatusInternalServerError)
		}

		return fiberCtx.SendStatus(fiber.StatusOK)
	}
}

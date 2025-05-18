package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/collection"
	"physk/internal/domain/aggregates/collection/entities"
	"physk/internal/domain/aggregates/user"
	uCtx "physk/internal/domain/aggregates/user/context"
	"physk/internal/infrastructure/db/storage"
	deliveryDTO "physk/internal/infrastructure/delivery/dto"
)

var (
	UserNotFound       = errors.New("user not found")
	CollectionNotFound = errors.New("collection not found")
	ImageNotFound      = errors.New("image not found")
)

func (c *Controller) RegisterUser(
	ctx context.Context,
	request deliveryDTO.UserCreateRequest,
) (user.User, error) {
	dto := request.ToAggregateDTO()
	userCtx := uCtx.NewUserContext(ctx)

	userAgg, err := user.CreateUser(userCtx, dto)
	if err != nil {
		return user.User{}, fmt.Errorf("create user: %w", err)
	}

	err = c.write.CreateUser(ctx, userAgg)
	if err != nil {
		return user.User{}, fmt.Errorf("create in write model: %w", err)
	}

	return userAgg, nil
}

func (c *Controller) Login(
	ctx context.Context,
	request deliveryDTO.UserLoginRequest,
) (string, error) {
	dto := request.ToAggregateDTO()
	userCtx := uCtx.NewUserContext(ctx)

	userAgg, err := c.read.GetUserByLogin(userCtx, dto.Login)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return "", UserNotFound
		}
		return "", fmt.Errorf("get user by login: %w", err)
	}

	err = userAgg.MatchPassword(request.Password)
	if err != nil {
		return "", fmt.Errorf("match password: %w", err)
	}

	token, err := c.tokenService.GenerateToken(userAgg.ID, userAgg.Role)
	if err != nil {
		return "", fmt.Errorf("gen new token: %w", err)
	}

	return token, nil
}

func (c *Controller) CreateCollection(
	ctx context.Context,
	request deliveryDTO.CreateCollectionRequest,
) (collection.Collection, error) {
	dto := request.ToAggregateDTO()

	collectionAgg := collection.NewCollection(dto)

	err := c.write.CreateCollection(ctx, collectionAgg)
	if err != nil {
		return collection.Collection{}, fmt.Errorf("create collection: %w", err)
	}

	return collectionAgg, nil
}

func (c *Controller) AddImageToCollection(
	ctx context.Context,
	request deliveryDTO.AddImageToCollectionRequest,
) (uuid.UUID, error) {
	dto := request.ToAggregateDTO()
	collectionAgg, err := c.read.GetCollectionByID(ctx, dto.CollectionID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return uuid.UUID{}, CollectionNotFound
		}
		return uuid.UUID{}, fmt.Errorf("get collection by id: %w", err)
	}

	imgEntity := entities.Image{
		ID:           uuid.New(),
		CollectionID: collectionAgg.ID,
		Name:         dto.Name,
		ContentType:  dto.ContentType,
		ImageData:    dto.Data,
	}

	err = collectionAgg.AddImage(imgEntity, request.IsPreview)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("collection aggregate: add image: %w", err)
	}

	err = c.write.AddImageToCollection(ctx, imgEntity, request.IsPreview)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("write model: add image to collection: %w", err)
	}
	return imgEntity.ID, nil
}

func (c *Controller) DeleteCollection(
	ctx context.Context,
	request deliveryDTO.DeleteCollectionRequest,
) error {
	dto := request.ToAggregateDTO()

	err := c.write.DeleteCollection(ctx, dto.CollectionID)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

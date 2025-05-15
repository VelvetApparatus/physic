package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/collection/entities"
	"physk/internal/domain/aggregates/user"
	"physk/internal/infrastructure/db/storage"
)

func (c *Controller) GetMe(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	usr, err := c.read.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return user.User{}, UserNotFound
		}
		return user.User{}, fmt.Errorf("get user by login: %w", err)
	}
	return usr, nil
}

func (c *Controller) GetImageByID(
	ctx context.Context,
	imageID uuid.UUID,
) (entities.Image, error) {
	img, err := c.read.GetImageByID(ctx, imageID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			return entities.Image{}, ImageNotFound
		}
		return entities.Image{}, fmt.Errorf("get image by id: %w", err)
	}

	return img, nil
}

func (c *Controller) GetImageIDsByCollectionID(
	ctx context.Context,
	collectionID uuid.UUID,
) ([]uuid.UUID, error) {
	ids, err := c.read.GetImageIDsByCollectionID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get image ids by collection id: %w", err)
	}
	return ids, nil
}

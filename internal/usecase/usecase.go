package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"physk/internal/domain/aggregates/collection"
	"physk/internal/domain/aggregates/collection/entities"
	"physk/internal/domain/aggregates/user"
	otlpStorage "physk/internal/infrastructure/db/storage"
	s3Storage "physk/internal/infrastructure/s3/storage"
)

type WriteModel interface {
	CreateUser(ctx context.Context, u user.User) error
	CreateCollection(ctx context.Context, c collection.Collection) error
	AddImageToCollection(ctx context.Context, c collection.Collection, i entities.Image) error
	DeleteCollection(ctx context.Context, collectionID uuid.UUID) error
}

type ReadModel interface {
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (user.User, error)
	GetCollectionByID(ctx context.Context, collectionID uuid.UUID) (collection.Collection, error)
	GetImageIDsByCollectionID(ctx context.Context, collectionID uuid.UUID) ([]uuid.UUID, error)
	GetImageByID(ctx context.Context, imgID uuid.UUID) (entities.Image, error)
}

type UseCase struct {
	otlp otlpStorage.Storage
	s3   s3Storage.Storage
}

func NewUseCase(st otlpStorage.Storage, s3 s3Storage.Storage) *UseCase {
	return &UseCase{otlp: st, s3: s3}
}

func (u *UseCase) CreateUser(ctx context.Context, us user.User) error {
	return u.otlp.CreateUser(ctx, us)
}

func (u *UseCase) CreateCollection(ctx context.Context, c collection.Collection) error {
	return u.otlp.CreateCollection(ctx, c)
}

func (u *UseCase) AddImageToCollection(ctx context.Context, c collection.Collection, img entities.Image) error {
	err := u.otlp.AddImageToCollection(ctx, img)
	if err != nil {
		return fmt.Errorf("otlp: add image to collection: %w", err)
	}
	err = u.s3.SaveImageByName(ctx, entities.ImageName(img.CollectionID, img.ID), img.ImageData)
	if err != nil {
		return fmt.Errorf("save image by name: %w", err)
	}
	return nil
}

func (u *UseCase) DeleteCollection(ctx context.Context, collectionID uuid.UUID) error {
	err := u.otlp.DeleteCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	iIds, err := u.otlp.GetImageIDsByCollectionID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get image ids by collection id: %w", err)
	}

	var deleteErr error
	for _, id := range iIds {
		err = u.s3.DeleteImageByName(ctx, entities.ImageName(collectionID, id))
		if err != nil {
			deleteErr = errors.Join(err)
		}
	}

	if deleteErr != nil {
		slog.LogAttrs(
			ctx, slog.LevelError, "delete images",
			slog.String("error", deleteErr.Error()),
		)
	}

	return nil
}

func (u *UseCase) GetUserByLogin(ctx context.Context, login string) (user.User, error) {
	return u.otlp.GetUserByLogin(ctx, login)
}

func (u *UseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	return u.otlp.GetUserByID(ctx, userID)
}

func (u *UseCase) GetCollectionByID(ctx context.Context, collectionID uuid.UUID) (collection.Collection, error) {
	return u.otlp.GetCollectionByID(ctx, collectionID)
}

func (u *UseCase) GetImageIDsByCollectionID(ctx context.Context, collectionID uuid.UUID) ([]uuid.UUID, error) {
	return u.otlp.GetImageIDsByCollectionID(ctx, collectionID)
}

func (u *UseCase) GetImageByID(
	ctx context.Context,
	imgID uuid.UUID,
) (entities.Image, error) {
	img, err := u.otlp.GetImageByID(ctx, imgID)
	if err != nil {
		return entities.Image{}, fmt.Errorf("get image by id: %w", err)
	}

	imgRaw, err := u.s3.GetImageByName(ctx, entities.ImageName(img.CollectionID, imgID))
	if err != nil {
		return entities.Image{}, fmt.Errorf("s3: get image doc: %w", err)
	}

	img.ImageData = make([]byte, len(imgRaw))
	copy(img.ImageData, imgRaw)

	return img, nil
}

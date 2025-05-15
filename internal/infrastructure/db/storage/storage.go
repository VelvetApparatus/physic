package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"physk/internal/domain/aggregates/collection"
	"physk/internal/domain/aggregates/collection/entities"
	"physk/internal/domain/aggregates/user"
)

var (
	ErrorNotFound = errors.New("not found")
)

type Storage interface {
	CreateUser(ctx context.Context, user user.User) error
	GetUserByLogin(ctx context.Context, login string) (user.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (user.User, error)
	CreateCollection(ctx context.Context, c collection.Collection) error
	AddImageToCollection(ctx context.Context, c collection.Collection, i entities.Image) error
	GetCollectionByID(ctx context.Context, collectionID uuid.UUID) (collection.Collection, error)
	DeleteCollection(ctx context.Context, collectionID uuid.UUID) error
	GetImageIDsByCollectionID(ctx context.Context, collectionID uuid.UUID) ([]uuid.UUID, error)
	GetImageByID(ctx context.Context, imgID uuid.UUID) (entities.Image, error)
}

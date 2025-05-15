package storage

import (
	"context"
)

type Storage interface {
	DeleteImageByName(ctx context.Context, name string) error
	SaveImageByName(ctx context.Context, name string, raw []byte) error
	GetImageByName(ctx context.Context, name string) ([]byte, error)
}

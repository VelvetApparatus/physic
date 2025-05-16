package storage

import (
	"context"
	"errors"
)

var (
	ErrNoFile = errors.New("no such key")
)

type Storage interface {
	DeleteImageByName(ctx context.Context, name string) error
	SaveImageByName(ctx context.Context, name string, raw []byte, contentType string) error
	GetImageByName(ctx context.Context, name string) ([]byte, string, error)
}

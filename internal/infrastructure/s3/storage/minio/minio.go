package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	s3Storage "physk/internal/infrastructure/s3/storage"
)

const (
	bucketName = `images`
)

type mImpl struct {
	cli *minio.Client
}

func NewS3Storage(minioCli *minio.Client) s3Storage.Storage {
	return &mImpl{cli: minioCli}
}

func (m *mImpl) DeleteImageByName(ctx context.Context, name string) error {
	err := m.cli.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}
	return nil
}

func (m *mImpl) SaveImageByName(ctx context.Context, name string, raw []byte, contentType string) error {
	_, err := m.cli.PutObject(ctx, bucketName, name, bytes.NewReader(raw), int64(len(raw)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

func (m *mImpl) GetImageByName(ctx context.Context, name string) ([]byte, string, error) {

	obj, err := m.cli.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, "", s3Storage.ErrNoFile

		}
		return nil, "", fmt.Errorf("get object: %w", err)
	}

	stats, err := obj.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("object: get stats: %w", err)
	}

	rawFile, err := io.ReadAll(obj)
	if err != nil {
		return nil, "", fmt.Errorf("read obj: %w", err)
	}

	return rawFile, stats.ContentType, nil
}

var _ s3Storage.Storage = (*mImpl)(nil)

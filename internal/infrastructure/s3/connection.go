package s3

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"physk/internal/config"
)

func NewMinioClient(ctx context.Context) (*minio.Client, error) {
	conf := config.C().S3

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: conf.UseSSL,
	}

	return minio.New(conf.Host+":"+conf.Port, opts)
}

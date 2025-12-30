package minioStorage

import (
	"context"
	"devconnectstorage/internal/domain"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(
	endpoint string,
	accessKey string,
	secretKey string,
	useSSL bool,
	bucket string,
) (*MinIOStorage, error) {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOStorage{
		client: client,
		bucket: bucket,
	}, nil
}

func (storage MinIOStorage) SaveFile(ctx context.Context, fileBytes io.Reader, file domain.File) (string, error) {
	objectName := buildObjectKey(file)
	info, err := storage.client.PutObject(ctx, storage.bucket, objectName, fileBytes, file.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func (storage MinIOStorage) DeleteFile(ctx context.Context, file domain.File) error {
	if file.StorageKey() == "" {
		return errors.New("storage key nil on delete")
	}

	return storage.client.RemoveObject(
		ctx,
		storage.bucket,
		file.StorageKey(),
		minio.RemoveObjectOptions{},
	)
}

func buildObjectKey(file domain.File) string {
	random := uuid.New()
	if file.ProjectID() != nil {
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			file.OwnerID(),
			random,
			*file.ProjectID(),
			file.FileName(),
		)
	}
	return fmt.Sprintf(
		"%s/%s/%s",
		file.OwnerID(),
		random,
		file.FileName(),
	)
}

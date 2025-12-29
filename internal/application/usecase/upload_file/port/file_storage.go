package port

import (
	"context"
	"devconnectstorage/internal/domain"
	"io"
)

type Storage interface {
	SaveFile(ctx context.Context, fileBytes io.Reader, file domain.File) (string, error)
	DeleteFile(ctx context.Context, file domain.File) error
}

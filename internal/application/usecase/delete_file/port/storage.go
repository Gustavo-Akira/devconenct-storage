package port

import (
	"context"
	"devconnectstorage/internal/domain"
)

type Storage interface {
	DeleteFile(ctx context.Context, file domain.File) error
}

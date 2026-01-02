package port

import (
	"context"
	"devconnectstorage/internal/domain"
)

type Repository interface {
	DeleteFile(ctx context.Context, id string) error
	GetFile(ctx context.Context, id string) (domain.File, error)
}

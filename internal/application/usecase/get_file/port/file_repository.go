package port

import (
	"context"
	"devconnectstorage/internal/domain"
)

type FileRepository interface {
	GetFile(ctx context.Context, id string) (domain.File, error)
}

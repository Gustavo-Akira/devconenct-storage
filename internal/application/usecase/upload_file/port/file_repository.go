package port

import (
	"context"
	"devconnectstorage/internal/domain"
)

type FileRepository interface {
	Save(ctx context.Context, file domain.File) (domain.File, error)
}

package getfile

import (
	"context"
	"devconnectstorage/internal/application/aggregate"
)

type IGetFileByIdUseCase interface {
	Execute(ctx context.Context, query GetFileByIdQuery) (*aggregate.FileContent, error)
}

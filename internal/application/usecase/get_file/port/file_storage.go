package port

import (
	"context"
	"io"
)

type FileStorage interface {
	GetFile(ctx context.Context, storageKey string) (io.ReadCloser, error)
}

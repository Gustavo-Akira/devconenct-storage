package aggregate

import (
	"devconnectstorage/internal/domain"
	"io"
)

type FileContent struct {
	Metadata domain.File
	Content  io.ReadCloser
}

package port

import "io"

type FileStorage interface {
	GetFile(storageKey string) (io.ReadCloser, error)
}

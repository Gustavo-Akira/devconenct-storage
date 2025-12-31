package port

import "devconnectstorage/internal/domain"

type FileRepository interface {
	GetFile(id string) (domain.File, error)
}

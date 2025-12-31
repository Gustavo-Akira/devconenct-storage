package uploadfile

import (
	"context"
	"devconnectstorage/internal/domain"
)

type IUploadFileUseCase interface {
	Execute(ctx context.Context, saveCommand UploadFileCommand) (domain.File, error)
}

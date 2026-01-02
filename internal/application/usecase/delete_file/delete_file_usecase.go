package deletefile

import (
	"context"
	"devconnectstorage/internal/application/usecase/delete_file/port"
)

type DeleteFileUseCase struct {
	repository port.Repository
	storage    port.Storage
}

func NewDeleteFileUseCase(repo port.Repository, storage port.Storage) *DeleteFileUseCase {
	return &DeleteFileUseCase{
		repository: repo,
		storage:    storage,
	}
}

func (uc *DeleteFileUseCase) Execute(ctx context.Context, command DeleteFileCommand) error {
	existentFile, err := uc.repository.GetFile(ctx, command.Id)
	if err != nil {
		return err
	}

	err = uc.storage.DeleteFile(ctx, existentFile)
	if err != nil {
		return err
	}
	err = uc.repository.DeleteFile(ctx, command.Id)

	return err
}

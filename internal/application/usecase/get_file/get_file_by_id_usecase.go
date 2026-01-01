package getfile

import (
	"context"
	"devconnectstorage/internal/application/aggregate"
	"devconnectstorage/internal/application/usecase/get_file/port"
)

type GetFileByIdUseCase struct {
	repository port.FileRepository
	storage    port.FileStorage
}

func NewGetFileByIdUseCase(repository port.FileRepository, storage port.FileStorage) *GetFileByIdUseCase {
	return &GetFileByIdUseCase{
		repository: repository,
		storage:    storage,
	}
}

func (uc *GetFileByIdUseCase) Execute(ctx context.Context, query GetFileByIdQuery) (*aggregate.FileContent, error) {
	metadata, repositoryError := uc.repository.GetFile(ctx, query.Id)
	if repositoryError != nil {
		return &aggregate.FileContent{}, repositoryError
	}
	content, storageError := uc.storage.GetFile(ctx, metadata.StorageKey())
	if storageError != nil {
		return &aggregate.FileContent{}, storageError
	}

	return &aggregate.FileContent{
		Metadata: metadata,
		Content:  content,
	}, nil
}

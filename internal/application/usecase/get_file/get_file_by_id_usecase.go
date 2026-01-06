package getfile

import (
	"context"
	"devconnectstorage/internal/application/aggregate"
	"devconnectstorage/internal/application/usecase/get_file/port"
	"devconnectstorage/internal/domain"
	"devconnectstorage/internal/infraestructure/outbound/auth"
	"errors"
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
	token := ctx.Value(auth.AuthTokenKey)

	if token == nil {
		return &aggregate.FileContent{}, errors.New("token cannot be null")
	}

	metadata, repositoryError := uc.repository.GetFile(ctx, query.Id)
	if repositoryError != nil {
		return &aggregate.FileContent{}, repositoryError
	}
	if metadata.Visibility() == domain.VisibilityPrivate && metadata.OwnerID() != token.(string) {
		return &aggregate.FileContent{}, errors.New("unauthorized")
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

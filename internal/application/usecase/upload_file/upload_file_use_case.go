package uploadfile

import (
	"context"
	"devconnectstorage/internal/application/usecase/upload_file/port"
	"devconnectstorage/internal/domain"
	"devconnectstorage/internal/infraestructure/outbound/auth"
	"errors"
	"strconv"
)

type UploadFileUseCase struct {
	fileRepository port.FileRepository
	storage        port.Storage
	generator      port.IdGenerator
	authClient     auth.IAuthClient
}

func NewUploadFileUseCase(repo port.FileRepository, storage port.Storage, generator port.IdGenerator, authClient auth.IAuthClient) *UploadFileUseCase {
	return &UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      generator,
		authClient:     authClient,
	}
}

func (uc UploadFileUseCase) Execute(ctx context.Context, saveCommand UploadFileCommand) (domain.File, error) {

	token := ctx.Value(auth.AuthTokenKey)

	if token == nil {
		return domain.File{}, errors.New("token cannot be null")
	}

	profileId, authError := uc.authClient.GetProfile(token.(string))

	if authError != nil {
		return domain.File{}, authError
	}

	file, domainErr := domain.NewFile(uc.generator.Generate(), strconv.FormatInt(*profileId, 10), saveCommand.ProjectID, saveCommand.FileName, saveCommand.MimeType, saveCommand.Size, domain.Visibility(saveCommand.Visibility))
	if domainErr != nil {
		return domain.File{}, domainErr
	}

	storageKey, storageErr := uc.storage.SaveFile(ctx, saveCommand.Content, file)

	if storageErr != nil {
		return domain.File{}, storageErr
	}

	if err := file.MarkAsAvailable(storageKey); err != nil {
		return domain.File{}, err
	}

	file, saveError := uc.fileRepository.Save(ctx, file)

	if saveError != nil {
		deleteError := uc.storage.DeleteFile(ctx, file)
		if deleteError != nil {
			return domain.File{}, deleteError
		}
		return domain.File{}, saveError
	}

	return file, nil
}

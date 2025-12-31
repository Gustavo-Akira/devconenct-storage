package uploadfile

import (
	"context"
	"devconnectstorage/internal/application/usecase/upload_file/port"
	"devconnectstorage/internal/domain"
)

type UploadFileUseCase struct {
	fileRepository port.FileRepository
	storage        port.Storage
	generator      port.IdGenerator
}

func NewUploadFileUseCase(repo port.FileRepository, storage port.Storage, generator port.IdGenerator) *UploadFileUseCase {
	return &UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      generator,
	}
}

func (uc UploadFileUseCase) Execute(ctx context.Context, saveCommand UploadFileCommand) (domain.File, error) {

	file, domainErr := domain.NewFile(uc.generator.Generate(), saveCommand.OwnerID, saveCommand.ProjectID, saveCommand.FileName, saveCommand.MimeType, saveCommand.Size, domain.Visibility(saveCommand.Visibility))
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

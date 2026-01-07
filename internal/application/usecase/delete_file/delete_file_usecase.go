package deletefile

import (
	"context"
	"devconnectstorage/internal/application/usecase/delete_file/port"
	"devconnectstorage/internal/infraestructure/outbound/auth"
	"errors"
	"strconv"
)

type DeleteFileUseCase struct {
	repository port.Repository
	storage    port.Storage
	authClient auth.IAuthClient
}

func NewDeleteFileUseCase(repo port.Repository, storage port.Storage, authClient auth.IAuthClient) *DeleteFileUseCase {
	return &DeleteFileUseCase{
		repository: repo,
		storage:    storage,
		authClient: authClient,
	}
}

func (uc *DeleteFileUseCase) Execute(ctx context.Context, command DeleteFileCommand) error {

	token := ctx.Value(auth.AuthTokenKey)

	if token == nil {
		return errors.New("token cannot be null")
	}

	profileId, authError := uc.authClient.GetProfile(token.(string))

	if authError != nil {
		return authError
	}

	existentFile, err := uc.repository.GetFile(ctx, command.Id)
	if err != nil {
		return err
	}

	if existentFile.OwnerID() != strconv.FormatInt(*profileId, 10) {
		return errors.New("unauthorized")
	}

	err = uc.storage.DeleteFile(ctx, existentFile)
	if err != nil {
		return err
	}
	err = uc.repository.DeleteFile(ctx, command.Id)

	return err
}

package deletefile

import (
	"context"
	"errors"
	"testing"

	"devconnectstorage/internal/application/usecase/delete_file/port"
	"devconnectstorage/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) GetFile(ctx context.Context, id string) (domain.File, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.File), args.Error(1)
}

func (m *RepositoryMock) DeleteFile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

var _ port.Repository = (*RepositoryMock)(nil)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) DeleteFile(ctx context.Context, file domain.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

var _ port.Storage = (*StorageMock)(nil)

func TestDeleteFileUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	repo := new(RepositoryMock)
	storage := new(StorageMock)

	file, _ := domain.NewFile(
		"1",
		"1fds",
		nil,
		"fsfdsa",
		"plain/text",
		32,
		domain.VisibilityPrivate,
	)

	repo.
		On("GetFile", ctx, "file-id").
		Return(file, nil)

	storage.
		On("DeleteFile", ctx, file).
		Return(nil)

	repo.
		On("DeleteFile", ctx, "file-id").
		Return(nil)

	uc := DeleteFileUseCase{
		repository: repo,
		storage:    storage,
	}

	err := uc.Execute(ctx, DeleteFileCommand{Id: "file-id"})

	assert.NoError(t, err)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestDeleteFileUseCase_Execute_ErrorOnGetFile(t *testing.T) {
	ctx := context.Background()

	repo := new(RepositoryMock)
	storage := new(StorageMock)

	expectedErr := errors.New("not found")

	repo.
		On("GetFile", ctx, "file-id").
		Return(domain.File{}, expectedErr)

	uc := DeleteFileUseCase{
		repository: repo,
		storage:    storage,
	}

	err := uc.Execute(ctx, DeleteFileCommand{Id: "file-id"})

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
	storage.AssertNotCalled(t, "DeleteFile")
}

func TestDeleteFileUseCase_Execute_ErrorOnStorageDelete(t *testing.T) {
	ctx := context.Background()

	repo := new(RepositoryMock)
	storage := new(StorageMock)

	file, _ := domain.NewFile(
		"1",
		"1fds",
		nil,
		"fsfdsa",
		"plain/text",
		32,
		domain.VisibilityPrivate,
	)
	expectedErr := errors.New("storage error")

	repo.
		On("GetFile", ctx, "file-id").
		Return(file, nil)

	storage.
		On("DeleteFile", ctx, file).
		Return(expectedErr)

	uc := DeleteFileUseCase{
		repository: repo,
		storage:    storage,
	}

	err := uc.Execute(ctx, DeleteFileCommand{Id: "file-id"})

	assert.Equal(t, expectedErr, err)

	repo.AssertNotCalled(t, "DeleteFile", ctx, "file-id")
}

func TestDeleteFileUseCase_Execute_ErrorOnRepositoryDelete(t *testing.T) {
	ctx := context.Background()

	repo := new(RepositoryMock)
	storage := new(StorageMock)

	file, _ := domain.NewFile(
		"1",
		"1fds",
		nil,
		"fsfdsa",
		"plain/text",
		32,
		domain.VisibilityPrivate,
	)
	expectedErr := errors.New("db error")

	repo.
		On("GetFile", ctx, "file-id").
		Return(file, nil)

	storage.
		On("DeleteFile", ctx, file).
		Return(nil)

	repo.
		On("DeleteFile", ctx, "file-id").
		Return(expectedErr)

	uc := DeleteFileUseCase{
		repository: repo,
		storage:    storage,
	}

	err := uc.Execute(ctx, DeleteFileCommand{Id: "file-id"})

	assert.Equal(t, expectedErr, err)
}

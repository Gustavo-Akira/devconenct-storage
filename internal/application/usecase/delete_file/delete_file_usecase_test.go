package deletefile

import (
	"context"
	"errors"
	"testing"

	"devconnectstorage/internal/domain"
	"devconnectstorage/internal/infraestructure/outbound/auth"

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

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) DeleteFile(ctx context.Context, file domain.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

type AuthClientMock struct {
	mock.Mock
}

func (m *AuthClientMock) GetProfile(token string) (*int64, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int64), args.Error(1)
}

func TestDeleteFileUseCase_Execute(t *testing.T) {
	const validToken = "valid-token"

	ctxWithToken := context.WithValue(context.Background(), auth.AuthTokenKey, validToken)

	setup := func() (*RepositoryMock, *StorageMock, *AuthClientMock, *DeleteFileUseCase) {
		repo := new(RepositoryMock)
		storage := new(StorageMock)
		authCli := new(AuthClientMock)
		uc := NewDeleteFileUseCase(repo, storage, authCli)
		return repo, storage, authCli, uc
	}

	t.Run("Success", func(t *testing.T) {
		repo, storage, authCli, uc := setup()
		ownerIDStr := "123"
		var ownerIDInt int64 = 123
		file, _ := domain.NewFile("1", ownerIDStr, nil, "path", "text/plain", 32, domain.VisibilityPrivate)

		authCli.On("GetProfile", validToken).Return(&ownerIDInt, nil)
		repo.On("GetFile", ctxWithToken, "file-id").Return(file, nil)
		storage.On("DeleteFile", ctxWithToken, file).Return(nil)
		repo.On("DeleteFile", ctxWithToken, "file-id").Return(nil)

		err := uc.Execute(ctxWithToken, DeleteFileCommand{Id: "file-id"})

		assert.NoError(t, err)
		authCli.AssertExpectations(t)
		repo.AssertExpectations(t)
		storage.AssertExpectations(t)
	})

	t.Run("Error No Token In Context", func(t *testing.T) {
		_, _, _, uc := setup()
		err := uc.Execute(context.Background(), DeleteFileCommand{Id: "file-id"})

		assert.EqualError(t, err, "token cannot be null")
	})

	t.Run("Error Auth Client GetProfile", func(t *testing.T) {
		repo, _, authCli, uc := setup()
		expectedErr := errors.New("auth error")

		authCli.On("GetProfile", validToken).Return(nil, expectedErr)

		err := uc.Execute(ctxWithToken, DeleteFileCommand{Id: "file-id"})

		assert.Equal(t, expectedErr, err)
		repo.AssertNotCalled(t, "GetFile", mock.Anything, mock.Anything)
	})

	t.Run("Error Unauthorized Owner", func(t *testing.T) {
		repo, storage, authCli, uc := setup()
		differentOwnerID := int64(456)
		file, _ := domain.NewFile("1", "999", nil, "path", "text/plain", 32, domain.VisibilityPrivate)

		authCli.On("GetProfile", validToken).Return(&differentOwnerID, nil)
		repo.On("GetFile", ctxWithToken, "file-id").Return(file, nil)

		err := uc.Execute(ctxWithToken, DeleteFileCommand{Id: "file-id"})

		assert.EqualError(t, err, "unauthorized")
		storage.AssertNotCalled(t, "DeleteFile", mock.Anything, mock.Anything)
	})

	t.Run("Error On GetFile", func(t *testing.T) {
		repo, _, authCli, uc := setup()
		expectedErr := errors.New("not found")
		var ownerIDInt int64 = 123
		authCli.On("GetProfile", validToken).Return(&ownerIDInt, nil)
		repo.On("GetFile", ctxWithToken, "file-id").Return(domain.File{}, expectedErr)

		err := uc.Execute(ctxWithToken, DeleteFileCommand{Id: "file-id"})

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Error On Storage Delete", func(t *testing.T) {
		repo, storage, authCli, uc := setup()
		ownerIDStr := "123"
		var ownerIDInt int64 = 123
		file, _ := domain.NewFile("1", ownerIDStr, nil, "path", "text/plain", 32, domain.VisibilityPrivate)
		expectedErr := errors.New("storage error")

		authCli.On("GetProfile", validToken).Return(&ownerIDInt, nil)
		repo.On("GetFile", ctxWithToken, "file-id").Return(file, nil)
		storage.On("DeleteFile", ctxWithToken, file).Return(expectedErr)

		err := uc.Execute(ctxWithToken, DeleteFileCommand{Id: "file-id"})

		assert.Equal(t, expectedErr, err)
		repo.AssertNotCalled(t, "DeleteFile", ctxWithToken, "file-id")
	})
}

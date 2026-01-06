package getfile

import (
	"bytes"
	"context"
	"devconnectstorage/internal/application/aggregate"
	"devconnectstorage/internal/domain"
	"devconnectstorage/internal/infraestructure/outbound/auth"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockRepositoryPort struct {
	mock func(ctx context.Context, id string) (domain.File, error)
}

func (repo *MockRepositoryPort) GetFile(ctx context.Context, id string) (domain.File, error) {
	return repo.mock(ctx, id)
}

type MockStoragePort struct {
	mock func(ctx context.Context, storageKey string) (io.ReadCloser, error)
}

func (storage *MockStoragePort) GetFile(ctx context.Context, storageKey string) (io.ReadCloser, error) {
	return storage.mock(ctx, storageKey)
}

func TestGetFileByIdUseCase_ShouldSuccess(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.AuthTokenKey, "sffa")
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("file content"))), nil
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.RehydrateFile(
				id,
				"owner 1",
				nil,
				"text.txt",
				"plain/text",
				12,
				"ssfdasdfdsfa",
				domain.VisibilityPublic,
				domain.StatusAvailable,
				time.Now(),
			)
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	returnedValue, err := usecase.Execute(ctx, query)
	require.NoError(t, err)
	assert.Equal(t, returnedValue.Metadata.ID(), "1")
	assert.NotEmpty(t, returnedValue.Content)
}

func TestGetFileByIdUseCase_ShouldSuccessWithPrivateVisibility(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.AuthTokenKey, "owner 1")
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("file content"))), nil
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.RehydrateFile(
				id,
				"owner 1",
				nil,
				"text.txt",
				"plain/text",
				12,
				"ssfdasdfdsfa",
				domain.VisibilityPrivate,
				domain.StatusAvailable,
				time.Now(),
			)
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	returnedValue, err := usecase.Execute(ctx, query)
	require.NoError(t, err)
	assert.Equal(t, returnedValue.Metadata.ID(), "1")
	assert.NotEmpty(t, returnedValue.Content)
}

func TestGetFileByIdUseCase_ShouldFailWithDifferentOwnerAndPrivate(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.AuthTokenKey, "sffa")
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("file content"))), nil
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.RehydrateFile(
				id,
				"owner 1",
				nil,
				"text.txt",
				"plain/text",
				12,
				"ssfdasdfdsfa",
				domain.VisibilityPrivate,
				domain.StatusAvailable,
				time.Now(),
			)
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	returnedValue, err := usecase.Execute(ctx, query)
	require.Error(t, err)
	assert.Equal(t, returnedValue, &aggregate.FileContent{})
}

func TestGetFileByIdUseCase_ShouldFailWithoutToken(t *testing.T) {
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("file content"))), nil
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.RehydrateFile(
				id,
				"owner 1",
				nil,
				"text.txt",
				"plain/text",
				12,
				"ssfdasdfdsfa",
				domain.VisibilityPublic,
				domain.StatusAvailable,
				time.Now(),
			)
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	_, err := usecase.Execute(context.Background(), query)
	assert.Error(t, err)
}

func TestGetFileByIdUseCase_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.AuthTokenKey, "sffa")
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte("file content"))), nil
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.File{}, errors.New("Some repo error")
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	_, err := usecase.Execute(ctx, query)
	require.Error(t, err)
}

func TestGetFileByIdUseCase_ShouldReturnErrorWhenStorageFails(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.AuthTokenKey, "sffa")
	mockStorage := MockStoragePort{
		mock: func(ctx context.Context, storageKey string) (io.ReadCloser, error) {
			return nil, errors.New("Error on storage")
		},
	}
	mockRepository := MockRepositoryPort{
		mock: func(ctx context.Context, id string) (domain.File, error) {
			return domain.RehydrateFile(
				id,
				"owner 1",
				nil,
				"text.txt",
				"plain/text",
				12,
				"ssfdasdfdsfa",
				domain.VisibilityPublic,
				domain.StatusAvailable,
				time.Now(),
			)
		},
	}

	usecase := NewGetFileByIdUseCase(
		&mockRepository,
		&mockStorage,
	)

	query := GetFileByIdQuery{
		Id: "1",
	}
	_, err := usecase.Execute(ctx, query)
	require.Error(t, err)
}

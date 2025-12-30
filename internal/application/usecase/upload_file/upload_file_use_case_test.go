package uploadfile

import (
	"bytes"
	"context"
	"devconnectstorage/internal/domain"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FileRepositoryMock struct {
	SaveFn func(ctx context.Context, file domain.File) (domain.File, error)
}

func (m *FileRepositoryMock) Save(ctx context.Context, file domain.File) (domain.File, error) {
	return m.SaveFn(ctx, file)
}

type FileStorageMock struct {
	SaveFileFn   func(ctx context.Context, content io.Reader, file domain.File) (string, error)
	DeleteFileFn func(ctx context.Context, file domain.File) error
}

func (m *FileStorageMock) SaveFile(ctx context.Context, content io.Reader, file domain.File) (string, error) {
	return m.SaveFileFn(ctx, content, file)
}

func (m *FileStorageMock) DeleteFile(ctx context.Context, file domain.File) error {
	return m.DeleteFileFn(ctx, file)
}

type IdGeneratorMock struct{}

func (gen *IdGeneratorMock) Generate() string {
	return "1"
}

func TestUploadFileUseCase_Success(t *testing.T) {
	ctx := context.Background()
	repoCalled := false
	storageCalled := false
	deleteStorageCalled := false

	repo := &FileRepositoryMock{
		SaveFn: func(ctx context.Context, file domain.File) (domain.File, error) {
			repoCalled = true
			return file, nil
		},
	}

	storage := &FileStorageMock{
		SaveFileFn: func(ctx context.Context, content io.Reader, file domain.File) (string, error) {
			storageCalled = true
			return "key1", nil
		},
		DeleteFileFn: func(ctx context.Context, file domain.File) error {
			deleteStorageCalled = true
			return nil
		},
	}

	uc := UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      &IdGeneratorMock{},
	}

	cmd := UploadFileCommand{
		OwnerID:    "owner-1",
		FileName:   "file.png",
		MimeType:   "image/png",
		Size:       123,
		Visibility: "PRIVATE",
		Content:    bytes.NewReader([]byte("file content")),
	}

	file, err := uc.Execute(ctx, cmd)

	assert.NoError(t, err)
	assert.True(t, repoCalled)
	assert.True(t, storageCalled)
	assert.False(t, deleteStorageCalled)
	assert.True(t, file.Status() == "AVAILABLE")
}

func TestUploadFileUseCase_Execute_ErrorOnRepositorySave(t *testing.T) {
	repoCalled := false
	storageCalled := false
	deleteStorageCalled := false
	ctx := context.Background()

	repo := &FileRepositoryMock{
		SaveFn: func(ctx context.Context, file domain.File) (domain.File, error) {
			return domain.File{}, errors.New("db error")
		},
	}

	storage := &FileStorageMock{
		SaveFileFn: func(ctx context.Context, content io.Reader, file domain.File) (string, error) {
			storageCalled = true
			return "key1", nil
		},
		DeleteFileFn: func(ctx context.Context, file domain.File) error {
			deleteStorageCalled = true
			return nil
		},
	}

	uc := UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      &IdGeneratorMock{},
	}

	cmd := UploadFileCommand{
		OwnerID:    "owner-id",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		Size:       1234,
		Visibility: string(domain.VisibilityPrivate),
		Content:    bytes.NewReader([]byte("file content")),
	}

	_, err := uc.Execute(ctx, cmd)

	assert.Error(t, err)
	assert.True(t, storageCalled)
	assert.False(t, repoCalled)
	assert.True(t, deleteStorageCalled)
}

func TestUploadFileUseCase_Execute_ErrorOnCreateFileDomain(t *testing.T) {
	repoCalled := false
	storageCalled := false
	deleteStorageCalled := false
	ctx := context.Background()

	repo := &FileRepositoryMock{
		SaveFn: func(ctx context.Context, file domain.File) (domain.File, error) {
			return domain.File{}, nil
		},
	}

	storage := &FileStorageMock{
		SaveFileFn: func(ctx context.Context, content io.Reader, file domain.File) (string, error) {
			storageCalled = true
			return "key1", nil
		},
		DeleteFileFn: func(ctx context.Context, file domain.File) error {
			deleteStorageCalled = true
			return nil
		},
	}

	uc := UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      &IdGeneratorMock{},
	}

	cmd := UploadFileCommand{
		OwnerID:    "",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		Size:       1234,
		Visibility: string(domain.VisibilityPrivate),
		Content:    bytes.NewReader([]byte("file content")),
	}

	_, err := uc.Execute(ctx, cmd)

	assert.Error(t, err)
	assert.False(t, storageCalled)
	assert.False(t, repoCalled)
	assert.False(t, deleteStorageCalled)
}

func TestUploadFileUseCase_Execute_ErrorOnMarkAvailableDomain(t *testing.T) {
	repoCalled := false
	storageCalled := false
	deleteStorageCalled := false
	ctx := context.Background()

	repo := &FileRepositoryMock{
		SaveFn: func(ctx context.Context, file domain.File) (domain.File, error) {
			return domain.File{}, nil
		},
	}

	storage := &FileStorageMock{
		SaveFileFn: func(ctx context.Context, content io.Reader, file domain.File) (string, error) {
			storageCalled = true
			return "", nil
		},
		DeleteFileFn: func(ctx context.Context, file domain.File) error {
			deleteStorageCalled = true
			return nil
		},
	}

	uc := UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      &IdGeneratorMock{},
	}

	cmd := UploadFileCommand{
		OwnerID:    "owner",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		Size:       1234,
		Visibility: string(domain.VisibilityPrivate),
		Content:    bytes.NewReader([]byte("file content")),
	}

	_, err := uc.Execute(ctx, cmd)

	assert.Error(t, err)
	assert.True(t, storageCalled)
	assert.False(t, repoCalled)
	assert.False(t, deleteStorageCalled)
}

func TestUploadFileUseCase_Execute_ErrorOnStorageSave(t *testing.T) {
	repoCalled := false
	storageCalled := false
	deleteStorageCalled := false
	ctx := context.Background()

	repo := &FileRepositoryMock{
		SaveFn: func(ctx context.Context, file domain.File) (domain.File, error) {
			storageCalled = true
			return file, nil
		},
	}

	storage := &FileStorageMock{
		SaveFileFn: func(ctx context.Context, content io.Reader, file domain.File) (string, error) {
			storageCalled = true
			return "", errors.New("Storage error")
		},
		DeleteFileFn: func(ctx context.Context, file domain.File) error {
			deleteStorageCalled = true
			return nil
		},
	}

	uc := UploadFileUseCase{
		fileRepository: repo,
		storage:        storage,
		generator:      &IdGeneratorMock{},
	}

	cmd := UploadFileCommand{
		OwnerID:    "owner-id",
		FileName:   "file.pdf",
		MimeType:   "application/pdf",
		Size:       1234,
		Visibility: string(domain.VisibilityPrivate),
		Content:    bytes.NewReader([]byte("file content")),
	}

	_, err := uc.Execute(ctx, cmd)

	assert.Error(t, err)
	assert.True(t, storageCalled)
	assert.False(t, repoCalled)
	assert.False(t, deleteStorageCalled)
}

func TestShouldCreateUseCaseWhenPassRepoAndStorage(t *testing.T) {
	fileRepositoryMock := &FileRepositoryMock{}
	fileStorageMock := &FileStorageMock{}
	generatorMock := &IdGeneratorMock{}
	uc := NewUploadFileUseCase(fileRepositoryMock, fileStorageMock, generatorMock)
	assert.Equal(t, uc.fileRepository, fileRepositoryMock)
	assert.Equal(t, uc.storage, fileStorageMock)
	assert.Equal(t, uc.generator, generatorMock)
}

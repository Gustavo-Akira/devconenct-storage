package minioStorage

import (
	"bytes"
	"context"
	"devconnectstorage/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinIOStorage_SaveAndDeleteFile_ShouldWork(t *testing.T) {
	ctx := context.Background()

	storage, err := NewMinIOStorage(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		false,
		"test-bucket",
	)
	require.NoError(t, err)

	content := []byte("file content")
	reader := bytes.NewReader(content)

	file, err := domain.NewFile(
		"owner-123",
		nil,
		"test.txt",
		"text/plain",
		int64(len(content)),
		domain.VisibilityPublic,
	)
	assert.NoError(t, err)

	key, err := storage.SaveFile(ctx, reader, file)
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	fileWithKey := file
	err = fileWithKey.MarkAsAvailable(key)
	require.NoError(t, err)

	err = storage.DeleteFile(ctx, fileWithKey)
	require.NoError(t, err)
}

func TestMinIOStorage_SaveAndDeleteFile_ShouldWorkWithProjectId(t *testing.T) {
	ctx := context.Background()

	storage, err := NewMinIOStorage(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		false,
		"test-bucket",
	)
	require.NoError(t, err)

	content := []byte("file content")
	reader := bytes.NewReader(content)
	projectId := "project122"
	file, err := domain.NewFile(
		"owner-123",
		&projectId,
		"test.txt",
		"text/plain",
		int64(len(content)),
		domain.VisibilityPublic,
	)
	assert.NoError(t, err)

	key, err := storage.SaveFile(ctx, reader, file)
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	fileWithKey := file
	err = fileWithKey.MarkAsAvailable(key)
	require.NoError(t, err)

	err = storage.DeleteFile(ctx, fileWithKey)
	require.NoError(t, err)
}

func TestMinIOStorage_SaveFile_ShouldFail_WhenBucketDoesNotExist(t *testing.T) {
	ctx := context.Background()

	storage, err := NewMinIOStorage(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		false,
		"bucket-not-exists",
	)
	require.NoError(t, err)

	reader := bytes.NewReader([]byte("content"))

	file, err := domain.NewFile(
		"owner",
		nil,
		"fail.txt",
		"text/plain",
		7,
		domain.VisibilityPublic,
	)
	require.NoError(t, err)

	key, err := storage.SaveFile(ctx, reader, file)

	assert.Error(t, err)
	assert.Empty(t, key)
}

func TestMinIOStorage_DeleteFile_ShouldBeNoOp_WhenStorageKeyIsEmpty(t *testing.T) {
	ctx := context.Background()

	storage, err := NewMinIOStorage(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		false,
		"test-bucket",
	)
	require.NoError(t, err)

	file, err := domain.NewFile(
		"owner",
		nil,
		"noop.txt",
		"text/plain",
		4,
		domain.VisibilityPublic,
	)
	require.NoError(t, err)

	err = storage.DeleteFile(ctx, file)

	assert.NoError(t, err)
}

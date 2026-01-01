package mongodb

import (
	"context"
	"testing"
	"time"

	"devconnectstorage/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMongoFileRepository_Save_ShouldPersistFile(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := db.Run(
		ctx,
		"mongo:8.2",
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = mongoContainer.Terminate(ctx)
	})

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := NewMongoFileRepository(
		mongoURI,
		"",
		"",
		"test-db",
		"files",
	)
	require.NoError(t, err)

	now := time.Now()
	file, err := domain.RehydrateFile(
		"file-123",
		"owner-123",
		nil,
		"test.txt",
		"text/plain",
		42,
		"storage/key",
		domain.VisibilityPublic,
		domain.StatusAvailable,
		now,
	)
	require.NoError(t, err)

	savedFile, err := repo.Save(ctx, file)
	require.NoError(t, err)

	assert.Equal(t, file.ID(), savedFile.ID())
	assert.Equal(t, file.OwnerID(), savedFile.OwnerID())
	assert.Equal(t, file.FileName(), savedFile.FileName())
	assert.Equal(t, file.Status(), savedFile.Status())
	var persisted domain.File
	persisted, err = repo.GetFile(ctx, file.ID())
	require.NoError(t, err)
	assert.Equal(t, file.ID(), persisted.ID())
	assert.Equal(t, file.OwnerID(), persisted.OwnerID())
	assert.Equal(t, file.FileName(), persisted.FileName())
	assert.Equal(t, file.Status(), persisted.Status())
	assert.Equal(t, file.Visibility(), persisted.Visibility())
	assert.WithinDuration(t, now, persisted.CreatedAt(), time.Second)
}

func TestMongoFileRepository_Save_ShouldNotPersistWithoutId(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := db.Run(
		ctx,
		"mongo:8.2",
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = mongoContainer.Terminate(ctx)
	})

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := NewMongoFileRepository(
		mongoURI,
		"",
		"",
		"",
		"files",
	)
	require.NoError(t, err)

	file, err := domain.NewFile(
		"2",
		"owner-123",
		nil,
		"test.txt",
		"text/plain",
		42,
		domain.VisibilityPublic,
	)
	require.NoError(t, err)

	_, err = repo.Save(ctx, file)
	assert.Error(t, err)

}

func TestMongoFileRepository_Get_ShouldReturnErrorWhenNotFind(t *testing.T) {
	ctx := context.Background()
	mongoContainer, err := db.Run(
		ctx,
		"mongo:8.2",
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = mongoContainer.Terminate(ctx)
	})

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := NewMongoFileRepository(
		mongoURI,
		"",
		"",
		"",
		"files",
	)
	require.NoError(t, err)

	_, err = repo.GetFile(ctx, "123")
	require.Error(t, err)
}

package minioStorage

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"devconnectstorage/internal/domain"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startMinioContainer(t *testing.T) (endpoint string, terminate func()) {
	req := tc.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForHTTP("/minio/health/ready").WithPort("9000/tcp").WithStartupTimeout(20 * time.Second),
		AutoRemove: true,
	}

	container, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(context.Background())
	require.NoError(t, err)
	port, err := container.MappedPort(context.Background(), "9000")
	require.NoError(t, err)

	return fmt.Sprintf("%s:%s", host, port.Port()), func() {
		_ = container.Terminate(context.Background())
	}
}

func createBucketForTest(t *testing.T, storage *MinIOStorage, bucketName string) {
	ctx := context.Background()
	exists, err := storage.client.BucketExists(ctx, bucketName)
	require.NoError(t, err)
	if !exists {
		err = storage.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		require.NoError(t, err)
	}
}

func TestMinIOStorage_WithTestContainer(t *testing.T) {
	endpoint, terminate := startMinioContainer(t)
	defer terminate()

	bucket := "test-bucket"

	client, err := NewMinIOStorage(endpoint, "minioadmin", "minioadmin", false, bucket)
	require.NoError(t, err)
	createBucketForTest(t, client, bucket)
	ctx := context.Background()
	content := []byte("file content")
	reader := bytes.NewReader(content)

	file, err := domain.NewFile("owner-123", nil, "test.txt", "text/plain", int64(len(content)), domain.VisibilityPublic)
	require.NoError(t, err)

	key, err := client.SaveFile(ctx, reader, file)
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	fileWithKey := file
	err = fileWithKey.MarkAsAvailable(key)
	require.NoError(t, err)

	err = client.DeleteFile(ctx, fileWithKey)
	require.NoError(t, err)
}

func TestMinIOStorage_ShoulReturnErrorOnDeleteWithoutStorageKey(t *testing.T) {
	ctx := context.Background()
	content := []byte("file content")
	file, err := domain.NewFile("owner-123", nil, "test.txt", "text/plain", int64(len(content)), domain.VisibilityPublic)
	require.NoError(t, err)
	client, err := NewMinIOStorage("localhost:9000", "minioadmin", "minioadmin", false, "test")
	require.NoError(t, err)
	err = client.DeleteFile(ctx, file)
	assert.Error(t, err)
}

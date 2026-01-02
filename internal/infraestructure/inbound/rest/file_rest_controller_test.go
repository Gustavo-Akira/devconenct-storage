package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"devconnectstorage/internal/application/aggregate"
	deletefile "devconnectstorage/internal/application/usecase/delete_file"
	getfile "devconnectstorage/internal/application/usecase/get_file"
	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"devconnectstorage/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UploadFileUseCaseMock struct {
	mock.Mock
}

type GetFileUseCaseMock struct {
	mock.Mock
}

type DeleteFileUseCaseMock struct {
	mock.Mock
}

func (m *GetFileUseCaseMock) Execute(ctx context.Context, query getfile.GetFileByIdQuery) (*aggregate.FileContent, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*aggregate.FileContent), args.Error(1)
}

func (m *UploadFileUseCaseMock) Execute(
	ctx context.Context,
	cmd uploadfile.UploadFileCommand,
) (domain.File, error) {
	args := m.Called(ctx, cmd)

	return args.Get(0).(domain.File), args.Error(1)
}

func (m *DeleteFileUseCaseMock) Execute(ctx context.Context, command deletefile.DeleteFileCommand) error {
	args := m.Called(ctx, command)
	return args.Error(0)
}

func TestUploadFile_ShouldReturn201_WhenRequestIsValid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(UploadFileUseCaseMock)
	controller := &FileRestController{
		uploadFile: useCaseMock,
	}

	router := gin.New()
	router.POST("/files", controller.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("owner_id", "owner-123")
	_ = writer.WriteField("visibility", "public")
	_ = writer.WriteField("project_id", "123")
	_ = writer.WriteField("mime_type", "plain/text")
	_ = writer.WriteField("file_name", "test.txt")

	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("file content"))

	_ = writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/files",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	expectedFile, _ := domain.NewFile(
		"123",
		"owner-123",
		nil,
		"test.txt",
		"text/plain",
		12,
		domain.VisibilityPublic,
	)

	useCaseMock.
		On("Execute", mock.Anything, mock.Anything).
		Return(expectedFile, nil).
		Once()

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestUploadFile_ShouldReturn400_WhenFileIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(UploadFileUseCaseMock)
	controller := &FileRestController{
		uploadFile: useCaseMock,
	}

	router := gin.New()
	router.POST("/files", controller.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("owner_id", "owner-123")
	_ = writer.WriteField("visibility", "public")
	_ = writer.WriteField("project_id", "123")
	_ = writer.WriteField("mime_type", "plain/text")
	_ = writer.WriteField("file_name", "test.txt")
	_ = writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/files",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	useCaseMock.AssertNotCalled(t, "Execute")
}

func TestUploadFile_ShouldReturn400_WhenRequiredMetadataIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(UploadFileUseCaseMock)
	controller := &FileRestController{
		uploadFile: useCaseMock,
	}

	router := gin.New()
	router.POST("/files", controller.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("ownerid", "owner-123")
	_ = writer.WriteField("visibility", "public")
	_ = writer.WriteField("project_id", "123")
	_ = writer.WriteField("mime_type", "plain/text")
	_ = writer.WriteField("file_name", "test.txt")
	_ = writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/files",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	useCaseMock.AssertNotCalled(t, "Execute")
}

func TestUploadFile_ShouldReturn400_WhenFileSizeIsZero(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(UploadFileUseCaseMock)
	controller := &FileRestController{uploadFile: useCaseMock}

	router := gin.New()
	router.POST("/files", controller.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("owner_id", "owner-123")
	_ = writer.WriteField("visibility", "public")
	_ = writer.WriteField("mime_type", "text/plain")
	_ = writer.WriteField("file_name", "empty.txt")

	_, _ = writer.CreateFormFile("file", "empty.txt")

	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	useCaseMock.AssertNotCalled(t, "Execute")
}

func TestUploadFile_ShouldReturn500_WhenUseCaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(UploadFileUseCaseMock)
	controller := &FileRestController{
		uploadFile: useCaseMock,
	}

	router := gin.New()
	router.POST("/files", controller.UploadFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("owner_id", "owner-123")
	_ = writer.WriteField("visibility", "public")
	_ = writer.WriteField("project_id", "123")
	_ = writer.WriteField("mime_type", "plain/text")
	_ = writer.WriteField("file_name", "test.txt")
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("content"))

	_ = writer.Close()

	req := httptest.NewRequest(
		http.MethodPost,
		"/files",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	useCaseMock.
		On("Execute", mock.Anything, mock.Anything).
		Return(domain.File{}, errors.New("use case error")).
		Once()

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestGetFileContentById_ShouldReturn200_WhenFileExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	content := bytes.NewBufferString("file content")
	file, _ := domain.NewFile(
		"123",
		"owner-1",
		nil,
		"test.txt",
		"text/plain",
		int64(content.Len()),
		domain.VisibilityPublic,
	)

	useCaseMock := new(GetFileUseCaseMock)
	controller := &FileRestController{
		getFile: useCaseMock,
	}

	router := gin.New()
	router.GET("/files/:id/content", controller.GetFileContentById)

	useCaseMock.On("Execute", mock.Anything, getfile.GetFileByIdQuery{Id: "123"}).
		Return(&aggregate.FileContent{Metadata: file, Content: io.NopCloser(content)}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/files/123/content", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "attachment; filename=\"test.txt\"", resp.Header().Get("Content-Disposition"))
	assert.Equal(t, "text/plain", resp.Header().Get("Content-Type"))
	assert.Equal(t, "file content", resp.Body.String())

	useCaseMock.AssertExpectations(t)
}

func TestGetFileContentById_ShouldReturn400_WhenIdMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &FileRestController{}
	router := gin.New()
	router.GET("/files/:id/content", controller.GetFileContentById)

	req := httptest.NewRequest(http.MethodGet, "/files//content", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetFileContentById_ShouldReturn500_WhenUseCaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(GetFileUseCaseMock)
	controller := &FileRestController{
		getFile: useCaseMock,
	}

	router := gin.New()
	router.GET("/files/:id/content", controller.GetFileContentById)

	useCaseMock.On("Execute", mock.Anything, mock.Anything).
		Return(&aggregate.FileContent{}, errors.New("use case error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/files/123/content", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestGetFileMetadataById_ShouldReturn200_WhenFileExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	content := bytes.NewBufferString("file content")
	file, _ := domain.NewFile(
		"123",
		"owner-1",
		nil,
		"test.txt",
		"text/plain",
		int64(content.Len()),
		domain.VisibilityPublic,
	)

	useCaseMock := new(GetFileUseCaseMock)
	controller := &FileRestController{
		getFile: useCaseMock,
	}

	router := gin.New()
	router.GET("/files/:id/metadata", controller.GetFileMetadataById)

	useCaseMock.On("Execute", mock.Anything, getfile.GetFileByIdQuery{Id: "123"}).
		Return(&aggregate.FileContent{Metadata: file, Content: io.NopCloser(content)}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/files/123/metadata", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), file.ID())
	assert.NotContains(t, resp.Body.String(), "storage_key")
	useCaseMock.AssertExpectations(t)
}

func TestGetFileMetadataById_ShouldReturn400_WhenIdMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	controller := &FileRestController{}
	router := gin.New()
	router.GET("/files/:id/metadata", controller.GetFileMetadataById)

	req := httptest.NewRequest(http.MethodGet, "/files//metadata", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetFileMetadataById_ShouldReturn500_WhenUseCaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(GetFileUseCaseMock)
	controller := &FileRestController{
		getFile: useCaseMock,
	}

	router := gin.New()
	router.GET("/files/:id/metadata", controller.GetFileMetadataById)

	useCaseMock.On("Execute", mock.Anything, mock.Anything).
		Return(&aggregate.FileContent{}, errors.New("use case error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/files/123/metadata", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestDeleteFile_ShouldReturn204WhenSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(DeleteFileUseCaseMock)
	controller := &FileRestController{
		deleteFile: useCaseMock,
	}

	router := gin.New()
	router.DELETE("/files/:id", controller.DeleteFile)

	useCaseMock.On("Execute", mock.Anything, mock.Anything).
		Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/files/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestDeleteFile_ShouldReturn500WhenFail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(DeleteFileUseCaseMock)
	controller := &FileRestController{
		deleteFile: useCaseMock,
	}

	router := gin.New()
	router.DELETE("/files/:id", controller.DeleteFile)

	useCaseMock.On("Execute", mock.Anything, mock.Anything).
		Return(errors.New("Error on delete")).Once()

	req := httptest.NewRequest(http.MethodDelete, "/files/123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	useCaseMock.AssertExpectations(t)
}

func TestDeleteFile_ShouldReturn400WhenIdMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useCaseMock := new(DeleteFileUseCaseMock)
	controller := &FileRestController{
		deleteFile: useCaseMock,
	}

	router := gin.New()
	router.DELETE("/files/:id/delete", controller.DeleteFile)

	req := httptest.NewRequest(http.MethodDelete, "/files//delete", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	useCaseMock.AssertExpectations(t)
}

package rest

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"devconnectstorage/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UploadFileUseCaseMock struct {
	mock.Mock
}

func (m *UploadFileUseCaseMock) Execute(
	ctx context.Context,
	cmd uploadfile.UploadFileCommand,
) (domain.File, error) {
	args := m.Called(ctx, cmd)

	return args.Get(0).(domain.File), args.Error(1)
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
	part.Write([]byte("file content"))

	writer.Close()

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
	writer.Close()

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
	writer.Close()

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

	writer.CreateFormFile("file", "empty.txt")

	writer.Close()

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
	part.Write([]byte("content"))

	writer.Close()

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

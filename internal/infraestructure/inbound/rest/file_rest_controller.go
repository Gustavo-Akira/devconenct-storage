package rest

import (
	"context"
	deletefile "devconnectstorage/internal/application/usecase/delete_file"
	getfile "devconnectstorage/internal/application/usecase/get_file"
	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"devconnectstorage/internal/infraestructure/inbound/rest/dto"
	"devconnectstorage/internal/infraestructure/outbound/auth"

	"github.com/gin-gonic/gin"
)

type FileRestController struct {
	uploadFile uploadfile.IUploadFileUseCase
	getFile    getfile.IGetFileByIdUseCase
	deleteFile deletefile.IDeleteFileUseCase
}

func NewFileRestController(usecase uploadfile.IUploadFileUseCase, getFileUsecase getfile.IGetFileByIdUseCase, deleteFileUseCase deletefile.IDeleteFileUseCase) *FileRestController {
	return &FileRestController{
		uploadFile: usecase,
		getFile:    getFileUsecase,
		deleteFile: deleteFileUseCase,
	}
}

func (controller *FileRestController) UploadFile(ctx *gin.Context) {
	var fileBody dto.UploadFileRequest
	jwt, err := ctx.Cookie("jwt")
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctxWithToken := context.WithValue(
		ctx.Request.Context(),
		auth.AuthTokenKey,
		jwt,
	)
	if err := ctx.ShouldBind(&fileBody); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "file is required"})
		return
	}

	if fileHeader.Size <= 0 {
		ctx.JSON(400, gin.H{"error": "file size cannot be 0"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		//govulncheck:ignore GO-2025-4233 reason: false positive via gin error handling; HTTP/3 not used
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			println(err)
		}
	}()

	command := fileBody.ToCommand(file, fileHeader.Size)

	result, err := controller.uploadFile.Execute(ctxWithToken, command)
	if err != nil {
		//govulncheck:ignore GO-2025-4233 reason: false positive via gin error handling; HTTP/3 not used
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, dto.NewFileMetadataResponse(result))
}

func (controller *FileRestController) GetFileContentById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "id cannot be empty"})
		return
	}

	jwt, err := ctx.Cookie("jwt")
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctxWithToken := context.WithValue(
		ctx.Request.Context(),
		auth.AuthTokenKey,
		jwt,
	)

	result, err := controller.getFile.Execute(
		ctxWithToken,
		getfile.GetFileByIdQuery{Id: id},
	)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer func() { _ = result.Content.Close() }()

	ctx.Header("Content-Disposition", "attachment; filename=\""+result.Metadata.FileName()+"\"")

	ctx.DataFromReader(
		200,
		result.Metadata.Size(),
		result.Metadata.MimeType(),
		result.Content,
		nil,
	)

}

func (controller *FileRestController) GetFileMetadataById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "id cannot be empty"})
		return
	}

	jwt, err := ctx.Cookie("jwt")
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctxWithToken := context.WithValue(
		ctx.Request.Context(),
		auth.AuthTokenKey,
		jwt,
	)

	result, err := controller.getFile.Execute(
		ctxWithToken,
		getfile.GetFileByIdQuery{Id: id},
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer func() { _ = result.Content.Close() }()
	ctx.JSON(200, dto.NewFileMetadataResponse(result.Metadata))
}

func (controller *FileRestController) DeleteFile(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "id cannot be empty"})
		return
	}
	command := deletefile.DeleteFileCommand{
		Id: id,
	}
	err := controller.deleteFile.Execute(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(204, gin.H{})
}

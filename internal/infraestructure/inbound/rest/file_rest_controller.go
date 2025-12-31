package rest

import (
	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"devconnectstorage/internal/infraestructure/inbound/rest/dto"

	"github.com/gin-gonic/gin"
)

type FileRestController struct {
	uploadFile uploadfile.IUploadFileUseCase
}

func NewFileRestController(usecase uploadfile.IUploadFileUseCase) *FileRestController {
	return &FileRestController{
		uploadFile: usecase,
	}
}

func (controller *FileRestController) UploadFile(ctx *gin.Context) {
	var fileBody dto.UploadFileRequest

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

	result, err := controller.uploadFile.Execute(ctx.Request.Context(), command)
	if err != nil {
		//govulncheck:ignore GO-2025-4233 reason: false positive via gin error handling; HTTP/3 not used
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, result)
}

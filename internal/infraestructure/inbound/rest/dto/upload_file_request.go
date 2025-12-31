package dto

import (
	uploadfile "devconnectstorage/internal/application/usecase/upload_file"
	"io"
)

type UploadFileRequest struct {
	OwnerID    string  `form:"owner_id" binding:"required"`
	ProjectID  *string `form:"project_id"`
	FileName   string  `form:"file_name" binding:"required"`
	MimeType   string  `form:"mime_type" binding:"required"`
	Visibility string  `form:"visibility" binding:"required"`
}

func (req UploadFileRequest) ToCommand(content io.Reader, size int64) uploadfile.UploadFileCommand {
	return uploadfile.UploadFileCommand{
		OwnerID:    req.OwnerID,
		ProjectID:  req.ProjectID,
		FileName:   req.FileName,
		MimeType:   req.MimeType,
		Size:       size,
		Visibility: req.Visibility,
		Content:    content,
	}
}

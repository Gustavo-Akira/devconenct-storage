package dto

import (
	"devconnectstorage/internal/domain"
	"time"
)

type FileMetadataResponse struct {
	Id         string    `json:"id"`
	OwnerID    string    `json:"owner_id"`
	ProjectID  *string   `json:"project_id"`
	FileName   string    `json:"file_name"`
	MimeType   string    `json:"mime_type"`
	Size       int64     `json:"size"`
	Visibility string    `json:"visibility"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewFileMetadataResponse(file domain.File) FileMetadataResponse {
	return FileMetadataResponse{
		Id:         file.ID(),
		OwnerID:    file.OwnerID(),
		ProjectID:  file.ProjectID(),
		FileName:   file.FileName(),
		MimeType:   file.MimeType(),
		Size:       file.Size(),
		Visibility: string(file.Visibility()),
		Status:     string(file.Status()),
		CreatedAt:  file.CreatedAt(),
	}
}

package mongodb

import (
	"devconnectstorage/internal/domain"
	"time"
)

type MongoFileEntity struct {
	ID         string    `bson:"_id"`
	OwnerID    string    `bson:"owner_id"`
	ProjectID  *string   `bson:"project_id,omitempty"`
	FileName   string    `bson:"file_name"`
	MimeType   string    `bson:"mime_type"`
	Size       int64     `bson:"size"`
	StorageKey string    `bson:"storage_key"`
	Visibility string    `bson:"visibility"`
	Status     string    `bson:"status"`
	CreatedAt  time.Time `bson:"created_at"`
}

func NewMongoFileEntity(file domain.File) MongoFileEntity {
	return MongoFileEntity{
		ID:         file.ID(),
		OwnerID:    file.OwnerID(),
		ProjectID:  file.ProjectID(),
		FileName:   file.FileName(),
		MimeType:   file.MimeType(),
		Size:       file.Size(),
		StorageKey: file.StorageKey(),
		Visibility: string(file.Visibility()),
		Status:     string(file.Status()),
		CreatedAt:  file.CreatedAt(),
	}
}

func (m *MongoFileEntity) ToDomain() (domain.File, error) {
	return domain.RehydrateFile(
		m.ID,
		m.OwnerID,
		m.ProjectID,
		m.FileName,
		m.MimeType,
		m.Size,
		m.StorageKey,
		domain.Visibility(m.Visibility),
		domain.Status(m.Status),
		m.CreatedAt,
	)
}

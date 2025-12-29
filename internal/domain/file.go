package domain

import (
	"fmt"
	"time"
)

type File struct {
	id         string
	ownerID    string
	projectID  *string
	fileName   string
	mimeType   string
	size       int64
	storageKey string
	visibility Visibility
	status     Status
	createdAt  time.Time
}

func createFile(id string, ownerID string, projectID *string, fileName string, mimeType string, size int64, storageKey string, visibility Visibility, status Status, createdAt time.Time) (File, error) {
	if ownerID == "" {
		return File{}, fmt.Errorf("ownerID cannot be empty")
	}
	if fileName == "" {
		return File{}, fmt.Errorf("fileName cannot be empty")
	}
	if size < 0 {
		return File{}, fmt.Errorf("size cannot be negative")
	}

	if visibility.IsValid() == false {
		return File{}, fmt.Errorf("invalid visibility value")
	}
	if status.IsValid() == false {
		return File{}, fmt.Errorf("invalid status value")
	}

	if createdAt.IsZero() {
		return File{}, fmt.Errorf("createdAt cannot be zero")
	}

	return File{
		id:         id,
		ownerID:    ownerID,
		projectID:  projectID,
		fileName:   fileName,
		mimeType:   mimeType,
		size:       size,
		storageKey: storageKey,
		visibility: visibility,
		status:     status,
		createdAt:  createdAt,
	}, nil
}

func NewFile(ownerID string, projectID *string, fileName string, mimeType string, size int64, visibility Visibility) (File, error) {
	return createFile("", ownerID, projectID, fileName, mimeType, size, "", visibility, StatusPending, time.Now())
}

func RehydrateFile(id string, ownerID string, projectID *string, fileName string, mimeType string, size int64, storageKey string, visibility Visibility, status Status, createdAt time.Time) (File, error) {
	if id == "" {
		return File{}, fmt.Errorf("id cannot be empty")
	}
	return createFile(id, ownerID, projectID, fileName, mimeType, size, storageKey, visibility, status, createdAt)
}

func (f File) ID() string {
	return f.id
}

func (f File) OwnerID() string {
	return f.ownerID
}

func (f File) ProjectID() *string {
	return f.projectID
}

func (f File) FileName() string {
	return f.fileName
}

func (f File) MimeType() string {
	return f.mimeType
}

func (f File) Size() int64 {
	return f.size
}

func (f File) StorageKey() string {
	return f.storageKey
}

func (f File) Visibility() Visibility {
	return f.visibility
}

func (f File) Status() Status {
	return f.status
}

func (f File) CreatedAt() time.Time {
	return f.createdAt
}

func (f *File) MarkAsAvailable(storageKey string) error {
	if f.status != StatusPending {
		return fmt.Errorf("file cannot be marked as available from %s", f.status)
	}
	if storageKey == "" {
		return fmt.Errorf("storageKey cannot be empty")
	}
	f.status = StatusAvailable
	f.storageKey = storageKey
	return nil
}

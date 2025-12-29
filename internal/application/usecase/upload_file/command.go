package uploadfile

import "io"

type UploadFileCommand struct {
	OwnerID    string
	ProjectID  *string
	FileName   string
	MimeType   string
	Size       int64
	Visibility string
	Content    io.Reader
}

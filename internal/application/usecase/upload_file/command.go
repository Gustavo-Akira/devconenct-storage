package uploadfile

import "io"

type UploadFileCommand struct {
	ProjectID  *string
	FileName   string
	MimeType   string
	Size       int64
	Visibility string
	Content    io.Reader
}

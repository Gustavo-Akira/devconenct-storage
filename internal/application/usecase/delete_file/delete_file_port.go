package deletefile

import "context"

type IDeleteFileUseCase interface {
	Execute(ctx context.Context, command DeleteFileCommand) error
}

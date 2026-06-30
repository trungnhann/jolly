package file

import "context"

type Storage interface {
	StoreFile(ctx context.Context, filePath string, fileContent []byte) (string, error)
}

package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type PublicStorage struct {
	rootDir string
	baseURL string
}

func NewPublicStorage(rootDir, baseURL string) *PublicStorage {
	if rootDir == "" {
		panic("rootDir cannot be empty")
	}
	if baseURL == "" {
		panic("baseURL cannot be empty")
	}

	return &PublicStorage{
		rootDir: rootDir,
		baseURL: baseURL,
	}
}

func (p *PublicStorage) StoreFile(ctx context.Context, filePath string, fileContent []byte) (string, error) {
	_ = ctx

	fullPath := filepath.Join(p.rootDir, filePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("create storage directory: %w", err)
	}

	if err := os.WriteFile(fullPath, fileContent, 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return p.baseURL + "/" + filepath.ToSlash(filePath), nil
}

package modelweights

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	SaveModelWeight(ctx context.Context, form *ModelUploadForm) (string, error)
}
type service struct {
	storageDir string
}

func NewService(storageDir string) Service {
	return &service{storageDir: storageDir}
}
func (s *service) SaveModelWeight(ctx context.Context, form *ModelUploadForm) (string, error) {
	ext := strings.ToLower(filepath.Ext(form.Weight.Filename))
	if ext != ".onnx" {
		return "", fmt.Errorf("invalid file extension %s, only .onnx is permitted", ext)
	}
	targetDir := filepath.Join(s.storageDir, form.UserId)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage matrix path: %w", err)
	}
	dstPath := filepath.Join(targetDir, "latest_model.onnx")

	src, err := form.Weight.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source stream: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open destination file descriptor: %w", err)
	}
	defer dst.Close()

	// 5. Stream from source to destination (safe memory consumption)
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to write blocks to disk: %w", err)
	}

	return dstPath, nil

}

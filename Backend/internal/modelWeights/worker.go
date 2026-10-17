package modelweights

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hibiken/asynq"
)

type ModelWeightWorker struct {
	storageDir string
}

func NewModelWeightWorker(storageDir string) *ModelWeightWorker {
	return &ModelWeightWorker{storageDir: storageDir}
}

func (w *ModelWeightWorker) HandleSaveModelWeight(ctx context.Context, t *asynq.Task) error {
	// 1. Unmarshal payload
	var payload SaveModelWeightPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to parse task payload: %w", err)
	}

	// CRITICAL: Always clean up the temporary file from the disk after execution
	defer os.Remove(payload.TempFilePath)

	// 2. Perform validations
	ext := strings.ToLower(filepath.Ext(payload.OriginalName))
	if ext != ".onnx" {
		return fmt.Errorf("invalid file extension %s, only .onnx is permitted", ext)
	}

	targetDir := filepath.Join(w.storageDir, payload.UserID)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	dstPath := filepath.Join(targetDir, "latest_model.onnx")

	// 3. Copy the file from the temporary directory to the final destination
	src, err := os.Open(payload.TempFilePath)
	if err != nil {
		return fmt.Errorf("failed to read temporary upload: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file to destination: %w", err)
	}

	log.Printf("Successfully saved model weights asynchronously for user: %s at path: %s", payload.UserID, dstPath)
	return nil
}

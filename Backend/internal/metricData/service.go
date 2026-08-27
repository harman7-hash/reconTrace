package metricData

import (
	"context"
	"dev/internal/prediction"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const TypeModelDownload = "model:download"

// Custom business errors
var ErrUnauthorized = errors.New("invalid or missing API key")

type ModelDownloadPayload struct {
	UserID string `json:"user_id"`
	S3Key  string `json:"s3_key"`
}

type Service interface {
	RecordMetric(ctx context.Context, apiKey string, metric *server_metric) error
	DailyCleanUp(ctx context.Context) error
	ProcessTrainingCompletion(ctx context.Context, payload TrainingCallbackPayload) error
}

type service struct {
	repo        Repository
	asynqClient *asynq.Client
	windows     *WindowStore
	batchWriter *BatchWriter
}

func NewService(repo Repository, asynqClient *asynq.Client, batchWriter *BatchWriter, windows *WindowStore) Service {
	return &service{
		repo:        repo,
		asynqClient: asynqClient,
		windows:     windows,
		batchWriter: batchWriter,
	}
}

// RecordMetric handles the primary business logic flow
func (s *service) RecordMetric(ctx context.Context, apiKey string, metric *server_metric) error {
	if apiKey == "" {
		return ErrUnauthorized
	}

	userID, err := s.repo.isValidAPIKey(ctx, apiKey)
	if err != nil {
		return err
	}

	if metric.Time.IsZero() {
		metric.Time = time.Now().UTC()
	}

	// if err := s.repo.InsertMetric(ctx, userID, metric); err != nil {
	// 	return err
	// }
	if err := s.batchWriter.Enqueue(ctx, userID, *metric); err != nil {
		return fmt.Errorf("failed to enqueue metric: %w", err)

	}
	s.windows.Push(userID, prediction.MetricPoint{
		CPUUtilization:    metric.CPUUtilization,
		MemoryUtilization: metric.MemoryUtilization,
		DiskReadBytes:     float64(metric.DiskReadBytes),
		DiskWriteBytes:    float64(metric.DiskWriteBytes),
	})
	return nil

}

func (s *service) DailyCleanUp(ctx context.Context) error {
	return s.repo.deleteMetrics(ctx)
}

// ProcessTrainingCompletion enqueues the model download task to Asynq
func (s *service) ProcessTrainingCompletion(ctx context.Context, payload TrainingCallbackPayload) error {
	if payload.Status == "FAILED" {
		log.Printf("[TRAINING FAILED] User %s: %s", payload.UserID, payload.ErrorMsg)
		return nil
	}

	taskPayload, err := json.Marshal(ModelDownloadPayload{
		UserID: payload.UserID,
		S3Key:  fmt.Sprintf("modelWeights/%s", payload.FileName),
	})
	fmt.Println(payload.FileName)
	if err != nil {
		return fmt.Errorf("failed to serialize task payload: %w", err)
	}

	task := asynq.NewTask(
		TypeModelDownload,
		taskPayload,
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Minute),
	)

	info, err := s.asynqClient.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue Asynq task: %w", err)
	}

	log.Printf("[ASYNQ ENQUEUED] Task ID: %s | Queue: %s | User: %s", info.ID, info.Queue, payload.UserID)
	return nil
}

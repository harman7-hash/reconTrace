package metricData

import (
	// "bytes"
	"bytes"
	"context"
	"dev/internal/tasks"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

const TypeMidnightPipeline = "cron:midnight_pipeline"

type Worker struct {
	db              *pgxpool.Pool
	httpClient      *http.Client
	asynqClient     *asynq.Client
	s3Client        *s3.Client
	s3BucketName    string
	modelStorageDir string
}

type UserTrainingPayload struct {
	UserID string `json:"user_id"`
}

func NewWorker(
	db *pgxpool.Pool,
	redisAddr string,
	s3Client *s3.Client,
	s3BucketName string,
	modelStorageDir string,
) *Worker {
	_ = os.MkdirAll(modelStorageDir, 0755)

	return &Worker{
		db:              db,
		httpClient:      &http.Client{Timeout: 30 * time.Second},
		asynqClient:     asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
		s3Client:        s3Client,
		s3BucketName:    s3BucketName,
		modelStorageDir: modelStorageDir,
	}
}

func (w *Worker) Close() error {
	return w.asynqClient.Close()
}

// func (w *Worker) HandleTrainingCoordinator(ctx context.Context, t *asynq.Task) error {
// 	rows, err := w.db.Query(ctx, `SELECT DISTINCT user_id FROM server_metrics`)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch users: %w", err)
// 	}
// 	defer rows.Close()

// 	var userIDs []string
// 	for rows.Next() {
// 		var id string
// 		if err := rows.Scan(&id); err == nil {
// 			userIDs = append(userIDs, id)
// 		}
// 	}

// 	for _, userID := range userIDs {
// 		payload := UserTrainingPayload{UserID: userID}
// 		jsonbytes, err := json.Marshal(payload)
// 		if err != nil {
// 			log.Printf("failed to marshal payload for user %s: %v", userID, err)
// 			continue
// 		}
// 		task := asynq.NewTask(tasks.TypeUserTraining, jsonbytes)

// 		info, err := w.asynqClient.Enqueue(task)
// 		if err != nil {
// 			log.Printf("failed to enqueue training task for user %s: %v", userID, err)
// 			continue
// 		}
// 		log.Printf("Enqueued training task for user %s (Task ID: %s)", userID, info.ID)
// 	}
// 	return nil
// }

func (w *Worker) StartTrainingScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		// Run check immediately on startup before waiting for first tick
		log.Println("[Scheduler] Running initial warmup/retrain check on startup...")
		w.processEligibleUsers(context.Background())

		for range ticker.C {
			log.Println("[Scheduler] Running hourly warmup/retrain check...")
			w.processEligibleUsers(context.Background())
		}
	}()
}

func (w *Worker) processEligibleUsers(ctx context.Context) {
	query := `
		SELECT id FROM users
		WHERE (model_status = 'WARMUP' AND NOW() >= registered_at + INTERVAL '14 days')
		   OR (model_status = 'READY'  AND NOW() >= last_trained_at + INTERVAL '14 days')
	`

	rows, err := w.db.Query(ctx, query)
	if err != nil {
		log.Printf("[Scheduler] Error querying eligible users: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			log.Printf("[Scheduler] Failed to scan user row: %v", err)
			continue
		}

		// 1. Lock user by setting status to TRAINING (prevents duplicate enqueues)
		_, updateErr := w.db.Exec(ctx, `UPDATE users SET model_status = 'TRAINING' WHERE id = $1`, userId)
		if updateErr != nil {
			log.Printf("[Scheduler] Failed to set model_status to TRAINING for user %s: %v", userId, updateErr)
			continue
		}

		// 2. Prepare task payload
		payload := UserTrainingPayload{UserID: userId}
		jsonbytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[Scheduler] Failed to marshal payload for user %s: %v", userId, err)
			continue
		}

		task := asynq.NewTask(tasks.TypeUserTraining, jsonbytes)

		// 3. Enqueue using Worker's existing asynqClient
		info, err := w.asynqClient.Enqueue(
			task,
			asynq.MaxRetry(3),
			asynq.Timeout(30*time.Minute),
			asynq.Queue("default"),
		)
		if err != nil {
			log.Printf("[Scheduler] Failed to enqueue training task for user %s: %v", userId, err)
			// Rollback status if Redis enqueue fails
			_, _ = w.db.Exec(ctx, `UPDATE users SET model_status = 'WARMUP' WHERE id = $1`, userId)
			continue
		}

		log.Printf("[Scheduler] Successfully enqueued training task for user %s (Task ID: %s)", userId, info.ID)
	}
}
func (w *Worker) HandleTraining(ctx context.Context, t *asynq.Task) error {
	var payload UserTrainingPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to parse user payload: %w", err)
	}
	log.Printf("Stage 2: Starting training for user: %s", payload.UserID)

	requestBody, err := json.Marshal(map[string]string{"user_id": payload.UserID})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://100.116.19.67:5000/train", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to dispatch training request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ML server returned error status: %s", resp.Status)
	}

	log.Printf("Stage 2: Successfully trained models for user: %s", payload.UserID)
	return nil
}

// HandleModelDownloadTask handles fetching model weights asynchronously
func (w *Worker) HandleModelDownloadTask(ctx context.Context, t *asynq.Task) error {
	var payload ModelDownloadPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w: %w", err, asynq.SkipRetry)
	}
	if err := w.fetchAndSaveFromS3(ctx, payload); err != nil {
		return fmt.Errorf("failed to download model from S3 for user %s: %w", payload.UserID, err)
		// fmt.Print("Hello")
	}

	return nil
}

// fetchAndSaveFromS3 downloads the file stream directly to disk
func (w *Worker) fetchAndSaveFromS3(ctx context.Context, payload ModelDownloadPayload) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(w.s3BucketName),
		Key:    aws.String(payload.S3Key),
	}

	result, err := w.s3Client.GetObject(ctx, input)
	if err != nil {
		return fmt.Errorf("S3 GetObject error for key %s: %w", payload.S3Key, err)
	}
	defer result.Body.Close()

	localPath := filepath.Join(w.modelStorageDir, fmt.Sprintf("%s", payload.UserID))
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, result.Body); err != nil {
		return fmt.Errorf("failed to write S3 stream to file: %w", err)
	}

	log.Printf("[S3 SUCCESS] Model stored via Asynq worker at: %s", localPath)
	return nil
}

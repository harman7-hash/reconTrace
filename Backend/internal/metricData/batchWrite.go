package metricData

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type queuedMetric struct {
	UserID string
	Metric server_metric
}

type BatchWriter struct {
	repo          Repository
	input         chan queuedMetric
	maxBatchSize  int
	flushInterval time.Duration

	closeOnce sync.Once
}

func NewBatchWriter(
	repo Repository,
	maxBatchSize int,
	flushInterval time.Duration,
) *BatchWriter {
	return &BatchWriter{
		repo:          repo,
		input:         make(chan queuedMetric, maxBatchSize*2),
		maxBatchSize:  maxBatchSize,
		flushInterval: flushInterval,
	}
}

// Enqueue accepts a metric for later database insertion.
// It waits if the in-memory queue is full, applying backpressure instead of
// silently losing metric data.
func (w *BatchWriter) Enqueue(
	ctx context.Context,
	userID string,
	metric server_metric,
) error {
	select {
	case w.input <- queuedMetric{
		UserID: userID,
		Metric: metric,
	}:
		return nil

	case <-ctx.Done():
		return fmt.Errorf("queue metric for batch insert: %w", ctx.Err())
	}
}

// Run must be started once in a goroutine.
// It flushes when either maxBatchSize is reached or flushInterval elapses.
func (w *BatchWriter) Run(ctx context.Context) {
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	batch := make([]queuedMetric, 0, w.maxBatchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}

		// Do not use the request context here: the HTTP request has usually
		// completed before this delayed database write runs.
		writeCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := w.repo.InsertMetricBatch(writeCtx, batch); err != nil {
			// Important: production code should retry this batch or place it in
			// a durable queue. Logging alone means a process crash/error can lose it.
			log.Printf("metric batch insert failed (%d records): %v", len(batch), err)
			return
		}

		log.Printf("stored metric batch: %d records", len(batch))
		batch = batch[:0]
	}

	for {
		select {
		case item := <-w.input:
			batch = append(batch, item)

			if len(batch) >= w.maxBatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-ctx.Done():
			// Best-effort final flush during graceful shutdown.
			flush()
			return
		}
	}
}

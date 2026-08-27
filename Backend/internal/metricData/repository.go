package metricData

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	isValidAPIKey(ctx context.Context, apikey string) (string, error)
	InsertMetricBatch(ctx context.Context, metrics_batch []queuedMetric) error
	deleteMetrics(ctx context.Context) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) isValidAPIKey(ctx context.Context, apiKey string) (string, error) {
	var userID string
	query := `SELECT id FROM users WHERE api_key = $1`

	err := r.db.QueryRow(ctx, query, apiKey).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to validate api key: %w", err)
	}

	return userID, nil
}

func (r *repository) InsertMetricBatch(ctx context.Context, metrics_batch []queuedMetric) error {
	for _, b := range metrics_batch {
		m := b.Metric
		query := `
			INSERT INTO server_metrics (
				time, user_id, cpu_util, mem_util, disk_r, disk_w
			) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := r.db.Exec(ctx, query,
			m.Time,
			b.UserID,
			m.CPUUtilization,
			m.MemoryUtilization,
			m.DiskReadBytes,
			m.DiskWriteBytes,
		)
		if err != nil {
			return fmt.Errorf("failed to insert metric hypertable row: %w", err)
		}
	}

	return nil
}
func (r *repository) deleteMetrics(ctx context.Context) error {
	query := `DELETE FROM server_metrics WHERE time < NOW() - INTERVAL '30 days'`
	_, err := r.db.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil

}

package prediction


import "time"

const WindowCapacity = 20

// MetricPoint represents a single time step of telemetry metrics.
type MetricPoint struct {
	CPUUtilization    float64 `json:"cpu_util"`
	MemoryUtilization float64 `json:"mem_util"`
	DiskReadBytes     float64 `json:"disk_r"`
	DiskWriteBytes    float64 `json:"disk_w"`
}

// PredictionInput represents a metric sequence window for anomaly evaluation.
type PredictionInput struct {
	Sequence  []MetricPoint `json:"sequence"`
	Threshold float64       `json:"threshold,omitempty"` // Custom threshold (falls back to default if 0)
	Timestamp time.Time     `json:"timestamp"`
}

// PredictionOutput contains the reconstructed MSE loss and anomaly decision.
type PredictionOutput struct {
	ReconstructionError float64   `json:"reconstruction_error"`
	Threshold           float64   `json:"threshold"`
	IsAnomaly           bool      `json:"is_anomaly"`
	PredictedAt         time.Time `json:"predicted_at"`
}

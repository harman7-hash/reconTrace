package prediction

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
)

type MetricWindowProvider interface {
	GetWindow(userID string) ([]MetricPoint, bool)
}
type Service interface {
	Predict(ctx context.Context, userID string, input PredictionInput) (PredictionOutput, error)
	UnloadModel(userID string) error
	Close() error
}

// type Service interface {
// 	Predict(ctx context.Context, userID string, input PredictionInput) (PredictionOutput, error)
// 	UnloadModel(userID string) error
// 	Close() error
// }

type service struct {
	storageDir  string
	windowStore MetricWindowProvider
	mu          sync.RWMutex
	predictors  map[string]Predictor
}

func NewService(storageDir string, windowStore MetricWindowProvider) Service {
	return &service{
		storageDir:  storageDir,
		windowStore: windowStore,
		predictors:  make(map[string]Predictor),
	}
}

func (s *service) Predict(ctx context.Context, userID string, input PredictionInput) (PredictionOutput, error) {
	if len(input.Sequence) == 0 {
		points, exists := s.windowStore.GetWindow(userID)
		if !exists || len(points) < WindowCapacity {
			actualCount := 0
			if exists {
				actualCount = len(points)
			}
			return PredictionOutput{}, fmt.Errorf("insufficient telemetry metrics collected: %d/%d", actualCount, WindowCapacity)
		}

		input.Sequence = points
	}
	predictor, err := s.getOrLoadPredictor(userID)
	if err != nil {
		return PredictionOutput{}, fmt.Errorf("failed to get predictor for user %s: %w", userID, err)
	}
	return predictor.Predict(input)
}

func (s *service) UnloadModel(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if predictor, exists := s.predictors[userID]; exists {
		delete(s.predictors, userID)
		return predictor.Close()
	}
	return nil
}

func (s *service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var firstErr error
	for userID, predictor := range s.predictors {
		if err := predictor.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(s.predictors, userID)
	}
	return firstErr
}

func (s *service) getOrLoadPredictor(userID string) (Predictor, error) {
	s.mu.RLock()
	predictor, exists := s.predictors[userID]
	s.mu.RUnlock()

	if exists {
		return predictor, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-checked locking pattern
	if predictor, exists = s.predictors[userID]; exists {
		return predictor, nil
	}

	modelPath := filepath.Join(s.storageDir, userID + ".onnx")
	fmt.Printf("%s", modelPath)
	newPredictor, err := NewONNXPredictor(modelPath)
	if err != nil {
		return nil, err
	}

	s.predictors[userID] = newPredictor
	return newPredictor, nil
}

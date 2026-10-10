package prediction

import (
	"fmt"
	"path/filepath"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

var (
	initOnce sync.Once
	initErr  error
)

func InitializeONNXRuntime(sharedLibPath string) error {
	initOnce.Do(func() {
		if ort.IsInitialized() {
			return
		}
		if sharedLibPath != "" {
			ort.SetSharedLibraryPath(sharedLibPath)
		}
		initErr = ort.InitializeEnvironment()
	})
	return initErr
}

type Predictor interface {
	Predict(input PredictionInput) (PredictionOutput, error)
	Close() error
}

type ONNXPredictor struct {
	session *ort.DynamicAdvancedSession
}

func NewONNXPredictor(modelPath string) (*ONNXPredictor, error) {
	if err := InitializeONNXRuntime(""); err != nil {
		return nil, fmt.Errorf("onnx runtime is not initialized: %w", err)
	}
	absPath, err := filepath.Abs(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	session, err := ort.NewDynamicAdvancedSession(
		absPath,
		[]string{"training_seq"},
		[]string{"cat"},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ONNX session: %w", err)
	}

	return &ONNXPredictor{session: session}, nil
}

func (p *ONNXPredictor) Predict(input PredictionInput) (PredictionOutput, error) {
	if len(input.Sequence) != WindowCapacity {
		return PredictionOutput{}, fmt.Errorf(
			"model requires exactly %d timestamps, received %d",
			WindowCapacity,
			len(input.Sequence),
		)
	}

	const featureDim = int64(4)
	timesteps := int64(WindowCapacity)
	totalElements := timesteps * featureDim

	inputData := make([]float32, 0, totalElements)
	for _, pt := range input.Sequence {
		inputData = append(inputData,
			float32(pt.CPUUtilization),
			float32(pt.MemoryUtilization),
			float32(pt.DiskReadBytes),
			float32(pt.DiskWriteBytes),
		)
	}

	inputShape := ort.NewShape(1, timesteps, featureDim)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return PredictionOutput{}, fmt.Errorf("failed to create input tensor: %w", err)
	}
	defer inputTensor.Destroy()

	outputData := make([]float32, totalElements)
	outputShape := ort.NewShape(1, timesteps, featureDim)
	outputTensor, err := ort.NewTensor(outputShape, outputData)
	if err != nil {
		return PredictionOutput{}, fmt.Errorf("failed to create output tensor: %w", err)
	}
	defer outputTensor.Destroy()

	if err := p.session.Run(
		[]ort.Value{inputTensor},
		[]ort.Value{outputTensor},
	); err != nil {
		return PredictionOutput{}, fmt.Errorf("failed to execute ONNX session: %w", err)
	}

	var sumSquaredDiff float64
	for i := range inputData {
		diff := float64(inputData[i] - outputData[i])
		sumSquaredDiff += diff * diff
	}

	mse := sumSquaredDiff / float64(totalElements)

	threshold := input.Threshold
	if threshold == 0 {
		threshold = 0.05
	}

	return PredictionOutput{
		ReconstructionError: mse,
		Threshold:           threshold,
		IsAnomaly:           mse > threshold,
		PredictedAt:         input.Timestamp,
	}, nil
}

func (p *ONNXPredictor) Close() error {
	if p.session != nil {
		return p.session.Destroy()
	}
	return nil
}

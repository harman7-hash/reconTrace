package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/metricData"
	modelweights "backend/internal/modelWeights"
	"backend/internal/tasks"
	"context"
	"fmt"
	"log"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hibiken/asynq"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cfg, err := config.Load()
	if err != nil {

	}
	defer cancel()
	// var databaseURI string
	dbPool, err := database.Connect(ctx, cfg.DatabaseURI)

	if err != nil {
		log.Fatalf("worker failed to connect DB: %v", err)

	}
	defer dbPool.Close()

	awsSDKCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				"",
			),
		),
	)
	s3Client := s3.NewFromConfig(awsSDKCfg)

	midnightWorker := metricData.NewWorker(
		dbPool,
		fmt.Sprintf("%s:6379", cfg.Host),
		s3Client,
		cfg.S3BucketName,
		"./weights",
	)
	modelWeightWorker := modelweights.NewModelWeightWorker("./weights")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:6379", cfg.Host)},
		asynq.Config{
			Concurrency: 10,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				if n >= 3 {
					return -1
				}
				return time.Duration(n*n) * time.Second * 10
			},
		},
	)

	mux := asynq.NewServeMux()

	// mux.HandleFunc(tasks.TypeTraining, midnightWorker.HandleTraining)
	// mux.HandleFunc(tasks.TypeUserTraining, midnightWorker.HandleTrainingCoordinator)
	mux.HandleFunc(metricData.TypeModelDownload, midnightWorker.HandleModelDownloadTask)
	mux.HandleFunc(tasks.TypeSaveModelWeight, modelWeightWorker.HandleSaveModelWeight)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Worker crash: %v", err)
	}

}

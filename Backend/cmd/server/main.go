package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/metricData"
	modelweights "backend/internal/modelWeights"
	"backend/internal/prediction"
	"backend/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func main() {
	// 1. Load Application Configurations
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading environment configurations: %v", err)
	}

	// 2. Database Migrations and Connection
	if err := database.RunMigrations(cfg.DatabaseURI); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to establish database connection pool: %v", err)
	}
	defer db.Close()

	// 3. Initialize Asynq Producer Client (For enqueuing tasks)
	redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:6379", cfg.Host)}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// 4. Initialize Repositories, Services, and Handlers
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	metricRepo := metricData.NewRepository(db)
	bgCtx, cancelBg := context.WithCancel(context.Background())
	defer cancelBg()

	batchWriter := metricData.NewBatchWriter(metricRepo, 100, 5*time.Second)
	go batchWriter.Run(bgCtx)

	sharedWindowStore := metricData.NewWindowStore()
	metricService := metricData.NewService(
		metricRepo,
		asynqClient,
		batchWriter,
		sharedWindowStore,
	)
	metricHandler := metricData.NewHandler(metricService)

	modelWeightsHandler := modelweights.NewHandler(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))

	predictionService := prediction.NewService("./weights", sharedWindowStore)
	defer func() {
		if err := predictionService.Close(); err != nil {
			log.Printf("close prediction service: %v", err)
		}
	}()
	predictionHandler := prediction.NewHTTPHandler(predictionService)

	// 5. Setup Gin Router and HTTP Routes
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	api := router.Group("/api")
	userHandler.RegisterRoutes(api)
	metricHandler.RegisterRoutes(api)
	modelWeightsHandler.RegisterRoutes(api)
	predictionHandler.RegisterRoutes(api)

	log.Printf("Server listening on http://%s:%s", cfg.Host, cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

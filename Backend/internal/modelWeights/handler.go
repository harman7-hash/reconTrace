package modelweights

import (
	"backend/internal/tasks"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type Handler struct {
	asynqClient *asynq.Client
}

type SaveModelWeightPayload struct {
	UserID       string `json:"user_id"`
	TempFilePath string `json:"temp_file_path"`
	OriginalName string `json:"original_name"`
}

func NewHandler(redisAddr string) *Handler {
	return &Handler{
		asynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/models/upload", h.HandleModelUpload)
}

func (h *Handler) HandleModelUpload(c *gin.Context) {
	// 1. Simple Security Header Token verification
	// if c.GetHeader("X-ML-Secret") != "Shizuka143" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized endpoint source"})
	// 	return
	// }

	var form ModelUploadForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing form fields: " + err.Error()})
		return
	}

	tempFile, err := os.CreateTemp("", "onnx_upload_*.tmp")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle file buffer: " + err.Error()})
		return
	}
	defer tempFile.Close()
	src, err := form.Weight.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open upload stream: " + err.Error()})
		return
	}
	defer src.Close()
	if _, err = io.Copy(tempFile, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write temp bytes: " + err.Error()})
		return
	}
	// 2. Prepare background task payload
	payload := SaveModelWeightPayload{
		UserID:       form.UserId,
		TempFilePath: tempFile.Name(),
		OriginalName: form.Weight.Filename,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payload error: " + err.Error()})
		return
	}
	// 3. Enqueue the task
	task := asynq.NewTask(tasks.TypeSaveModelWeight, jsonBytes)
	info, err := h.asynqClient.Enqueue(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue processing job: " + err.Error()})
		return
	}
	// 4. Return immediately with 202 Accepted (releases TCP connection)
	c.JSON(http.StatusAccepted, gin.H{
		"status":    "accepted",
		"message":   "Model weights uploaded successfully. Processing in background.",
		"task_id":   info.ID,
		"server_id": form.UserId,
	})
}

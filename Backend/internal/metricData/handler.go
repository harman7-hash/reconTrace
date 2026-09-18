package metricData

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

type TrainingCallbackPayload struct {
	UserID    string  `json:"user_id" binding:"required"`
	Status    string  `json:"status" binding:"required"` // "COMPLETED" or "FAILED"
	ModelURL  string  `json:"model_url,omitempty"`
	ErrorMsg  string  `json:"error_msg,omitempty"`
	Threshold float64 `json:"threshold" binding:"required"`
	FileName  string  `json:"filename" binding:"required"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes connects this handler to your main Gin router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/metrics", h.HandlePostMetric)
	router.POST("/training/callback", h.HandleTrainingCallback)
}

func (h *Handler) HandleTrainingCallback(c *gin.Context) {
	var payload TrainingCallbackPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid callback payload: " + err.Error()})
		return
	}
	if payload.Status == "FAILED" {
		log.Print("Training Failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Training at the ML server Failed"})
	}
	if err := h.service.ProcessTrainingCompletion(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue download task: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Callback processed successfully",
		"message": "Model weight retrieval queued for execution",
	})
}

func (h *Handler) HandlePostMetric(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	var metric server_metric
	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	
	err := h.service.RecordMetric(c.Request.Context(), apiKey, &metric)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while processing metric"})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{"status": "Metric recorded successfully"})
}

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
	// 1. Extract the API key from the HTTP Headers (Standard industry practice)
	apiKey := c.GetHeader("X-API-Key")

	// 2. Parse the incoming JSON body into our Entity struct
	var metric server_metric
	if err := c.ShouldBindJSON(&metric); err != nil {
		// ShouldBindJSON automatically catches missing required fields (like server_id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	// 3. Pass the extracted data down to the Service layer
	// Notice we pass c.Request.Context() so database queries can be cancelled if the user disconnects!
	err := h.service.RecordMetric(c.Request.Context(), apiKey, &metric)
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			return
		}
		// If it's a real database error, log it internally but return a generic 500 to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while processing metric"})
		return
	}

	// 4. Success! Return a 201 Created response.
	c.JSON(http.StatusCreated, gin.H{"status": "Metric recorded successfully"})
}

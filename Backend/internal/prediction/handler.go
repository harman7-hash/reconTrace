package prediction

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(svc Service) *HTTPHandler {
	return &HTTPHandler{service: svc}
}

func (h *HTTPHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users/:user_id/predict", h.Predict)
}

func (h *HTTPHandler) Predict(c *gin.Context) {
	userID := c.Param("user_id")
	fmt.Printf("%s", userID)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id route parameter is required"})
		return
	}

	var input PredictionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}

	if input.Timestamp.IsZero() {
		input.Timestamp = time.Now()
	}

	output, err := h.service.Predict(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

package api

import (
	"net/http"
	"time"

	"inci-backend/models"
	"inci-backend/observability"
	"inci-backend/worker"

	"github.com/gin-gonic/gin"
)

// IngestSignal handles the POST /api/v1/signals endpoint.
func IngestSignal(c *gin.Context) {
	var sig models.Signal
	if err := c.ShouldBindJSON(&sig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sig.Timestamp.IsZero() {
		sig.Timestamp = time.Now()
	}

	// Non-blocking push to the worker queue. 
	// Since buffer is 20000, we use a select with default to avoid blocking.
	select {
	case worker.JobQueue <- sig:
		observability.RecordSignal()
		// Return 202 immediately as requested, before DB writes occur
		c.Status(http.StatusAccepted)
	default:
		// Queue is completely full
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service overloaded"})
	}
}

// HealthCheck handles the /health endpoint.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
		"time":   time.Now(),
	})
}

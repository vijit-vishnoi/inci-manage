package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"inci-backend/db"
	"inci-backend/models"
	"inci-backend/workflow"

	retry "github.com/avast/retry-go/v4"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetIncidents handles GET /api/v1/incidents
func GetIncidents(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Read-Through Cache Check
	cachedData, err := db.RedisClient.Get(ctx, "active_incidents").Result()
	if err == nil && cachedData != "" {
		var items []models.WorkItem
		if err := json.Unmarshal([]byte(cachedData), &items); err == nil {
			c.JSON(http.StatusOK, items)
			return
		}
	}

	var items []models.WorkItem

	err = retry.Do(
		func() error {
			items = []models.WorkItem{} // Reset
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			query := `SELECT id, component_id, title, description, severity, status, created_at, updated_at 
			          FROM work_items WHERE status != 'CLOSED' ORDER BY severity ASC`
			rows, err := db.PGPool.Query(timeoutCtx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var w models.WorkItem
				if err := rows.Scan(&w.ID, &w.ComponentID, &w.Title, &w.Description, &w.Severity, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
					return err
				}
				items = append(items, w)
			}
			return rows.Err()
		},
		retry.Attempts(3),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Async Cache Save (5s TTL)
	go func(incidents []models.WorkItem) {
		bgCtx := context.Background()
		if data, err := json.Marshal(incidents); err == nil {
			db.RedisClient.Set(bgCtx, "active_incidents", data, 5*time.Second)
		}
	}(items)

	c.JSON(http.StatusOK, items)
}

// GetIncidentDetails handles GET /api/v1/incidents/:id
func GetIncidentDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var w models.WorkItem
	var rawSignals []bson.M

	err = retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 1. Get from Postgres
			query := `SELECT id, component_id, title, description, severity, status, created_at, updated_at 
			          FROM work_items WHERE id = $1`
			err := db.PGPool.QueryRow(ctx, query, id).Scan(
				&w.ID, &w.ComponentID, &w.Title, &w.Description, &w.Severity, &w.Status, &w.CreatedAt, &w.UpdatedAt,
			)
			if err != nil {
				return err
			}

			// 2. Get from MongoDB
			collection := db.MongoClient.Database("inci_mongo_db").Collection("raw_signals")
			cursor, err := collection.Find(ctx, bson.M{"component_id": w.ComponentID})
			if err != nil {
				return err
			}
			defer cursor.Close(ctx)
			
			rawSignals = []bson.M{}
			if err = cursor.All(ctx, &rawSignals); err != nil {
				return err
			}

			return nil
		},
		retry.Attempts(3),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incident":    w,
		"raw_signals": rawSignals,
	})
}

// SubmitRCA handles POST /api/v1/incidents/:id/rca
func SubmitRCA(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var rca models.RCA
	if err := c.ShouldBindJSON(&rca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rca.WorkItemID = id

	err = retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Calculate MTTR from explicit dates
			mttrMinutes := int(rca.IncidentEnd.Sub(rca.IncidentStart).Minutes())
			if mttrMinutes < 0 {
				mttrMinutes = 0 // Safeguard
			}
			rca.MTTRMinutes = mttrMinutes

			// 3. Insert RCA
			query := `INSERT INTO rca_records (work_item_id, root_cause_category, fix_applied, prevention_steps, mttr_minutes)
			          VALUES ($1, $2, $3, $4, $5)`
			_, err = db.PGPool.Exec(ctx, query, rca.WorkItemID, rca.RootCauseCategory, rca.FixApplied, rca.PreventionSteps, rca.MTTRMinutes)
			return err
		},
		retry.Attempts(3),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rca)
}

// UpdateStatus handles PATCH /api/v1/incidents/:id/status
func UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = retry.Do(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Get current status
			var currentStatus string
			err := db.PGPool.QueryRow(ctx, "SELECT status FROM work_items WHERE id = $1", id).Scan(&currentStatus)
			if err != nil {
				return err
			}

			// Validator for DB
			rcaValidator := func(ctx context.Context, wid int) error {
				var count int
				err := db.PGPool.QueryRow(ctx, "SELECT COUNT(*) FROM rca_records WHERE work_item_id = $1", wid).Scan(&count)
				if err != nil {
					return err
				}
				if count == 0 {
					return workflow.ErrIncompleteRCA
				}
				return nil
			}

			// Validate Transition
			if err := workflow.ValidateTransition(ctx, currentStatus, req.Status, id, rcaValidator); err != nil {
				return retry.Unrecoverable(err)
			}

			// Update Status
			_, err = db.PGPool.Exec(ctx, "UPDATE work_items SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", req.Status, id)
			return err
		},
		retry.Attempts(3),
	)

	if err != nil {
		if errors.Is(err, workflow.ErrIncompleteRCA) || errors.Is(err, workflow.ErrInvalidTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Cache Invalidation
	db.RedisClient.Del(context.Background(), "active_incidents")

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

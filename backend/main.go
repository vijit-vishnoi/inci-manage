package main

import (
	"context"
	"log"
	"time"

	"inci-backend/api"
	"inci-backend/db"
	"inci-backend/debouncer"
	"inci-backend/observability"
	"inci-backend/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize connections
	ctx := context.Background()
	// Using default localhost assuming it's running outside the docker network for local dev.
	// In production, these should be environment variables.
	pgConnString := "postgres://inci_user:inci_password@localhost:5432/inci_db"
	if err := db.InitPostgres(ctx, pgConnString); err != nil {
		log.Fatalf("Postgres init error: %v", err)
	}
	defer db.PGPool.Close()

	mongoURI := "mongodb://root:rootpassword@localhost:27017"
	if err := db.InitMongo(ctx, mongoURI); err != nil {
		log.Fatalf("Mongo init error: %v", err)
	}
	defer func() {
		if err := db.MongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Mongo disconnect error: %v", err)
		}
	}()

	redisAddr := "localhost:6379"
	if err := db.InitRedis(ctx, redisAddr); err != nil {
		log.Fatalf("Redis init error: %v", err)
	}
	defer db.RedisClient.Close()

	// Start Observability Logger (logs every 5 seconds)
	observability.StartLogger(5 * time.Second)

	// Initialize Debouncer (10 seconds window with in-memory map + sync.RWMutex)
	d := debouncer.New(10 * time.Second)

	// Initialize Worker Pool (100 workers, 20000 channel buffer size)
	worker.InitPool(100, 20000, d)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS Config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Middleware (12,000 req/sec limit to handle bursts)
	router.Use(api.RateLimiter(rate.Limit(12000), 12000))

	// Routes
	router.GET("/health", api.HealthCheck)
	
	v1 := router.Group("/api/v1")
	{
		v1.POST("/signals", api.IngestSignal)
		v1.GET("/incidents", api.GetIncidents)
		v1.GET("/incidents/:id", api.GetIncidentDetails)
		v1.POST("/incidents/:id/rca", api.SubmitRCA)
		v1.PATCH("/incidents/:id/status", api.UpdateStatus)
	}

	log.Println("Server starting on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

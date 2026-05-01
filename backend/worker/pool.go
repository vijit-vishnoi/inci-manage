package worker

import (
	"context"
	"fmt"
	"log"

	"inci-backend/db"
	"inci-backend/debouncer"
	"inci-backend/models"

	retry "github.com/avast/retry-go/v4"
)

var (
	JobQueue chan models.Signal
)

// InitPool initializes the buffered channel and starts the worker pool.
func InitPool(numWorkers int, bufferSize int, d *debouncer.Debouncer) {
	JobQueue = make(chan models.Signal, bufferSize)

	for i := 0; i < numWorkers; i++ {
		go worker(i, d)
	}
	log.Printf("Worker pool started with %d workers and buffer %d\n", numWorkers, bufferSize)
}

func worker(id int, d *debouncer.Debouncer) {
	for sig := range JobQueue {
		ctx := context.Background()

		// 1. Always insert into MongoDB raw_signals collection with retry
		err := retry.Do(
			func() error {
				collection := db.MongoClient.Database("inci_mongo_db").Collection("raw_signals")
				_, err := collection.InsertOne(ctx, sig)
				return err
			},
			retry.Attempts(3),
		)
		if err != nil {
			log.Printf("[Worker %d] Failed to insert raw signal to MongoDB: %v\n", id, err)
		}

		// 2. Debounce and conditionally insert to Postgres with retry
		if d.Allow(sig.ComponentID) {
			err = retry.Do(
				func() error {
					query := `
						INSERT INTO work_items (component_id, title, description, severity, status) 
						VALUES ($1, $2, $3, $4, 'OPEN')
					`
					title := fmt.Sprintf("Incident for Component %s", sig.ComponentID)
					desc := fmt.Sprintf("Error Code: %s", sig.ErrorCode)

					// Determine severity (default 3, or extract from metadata if exists)
					severity := 3
					if s, ok := sig.Metadata["severity"].(float64); ok {
						severity = int(s)
					}

					_, err := db.PGPool.Exec(ctx, query, sig.ComponentID, title, desc, severity)
					return err
				},
				retry.Attempts(3),
			)
			if err != nil {
				log.Printf("[Worker %d] Failed to insert work_item to Postgres: %v\n", id, err)
			}
		}
	}
}

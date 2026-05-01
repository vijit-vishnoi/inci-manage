package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PGPool      *pgxpool.Pool
	MongoClient *mongo.Client
	RedisClient *redis.Client
)

// InitPostgres connects to PostgreSQL and returns a pool.
func InitPostgres(ctx context.Context, connString string) error {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return fmt.Errorf("unable to create pg connection pool: %v", err)
	}
	PGPool = pool
	return nil
}

// InitMongo connects to MongoDB and returns a client.
func InitMongo(ctx context.Context, uri string) error {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("unable to connect to mongo: %v", err)
	}
	
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("mongo ping failed: %v", err)
	}
	MongoClient = client
	return nil
}

// InitRedis connects to Redis and returns a client.
func InitRedis(ctx context.Context, addr string) error {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("unable to connect to redis: %v", err)
	}
	RedisClient = client
	return nil
}

package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	URI        string
	Database   string
	Collection string
}

// LoadMongoConfig builds Mongo connection settings from env vars.
func LoadMongoConfig() (*MongoConfig, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URI must be set when DB_BACKEND=mongo")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		return nil, fmt.Errorf("MONGO_DB must be set when DB_BACKEND=mongo")
	}

	collection := os.Getenv("MONGO_COLLECTION")
	if collection == "" {
		collection = "users"
	}

	return &MongoConfig{
		URI:        uri,
		Database:   dbName,
		Collection: collection,
	}, nil
}

// NewMongoClient connects to MongoDB using the current environment config.
func NewMongoClient(ctx context.Context) (*mongo.Client, *MongoConfig, error) {
	cfg, err := LoadMongoConfig()
	if err != nil {
		return nil, nil, err
	}

	opts := options.Client().ApplyURI(cfg.URI)
	opts.SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// Ping to ensure connectivity before returning the client.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	slog.Info("connected to mongodb", slog.String("database", cfg.Database))

	return client, cfg, nil
}

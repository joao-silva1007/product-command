package core

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabase() (*mongo.Database, error) {
	databaseConnection := os.Getenv("MONGO_URL")
	databaseName := os.Getenv("DB_NAME")

	client, err := mongo.NewClient(options.Client().ApplyURI(databaseConnection))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	database := client.Database(databaseName)

	if err != nil {
		return nil, err
	}

	return database, nil
}

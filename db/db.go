package db

import (
	"context"
	"fmt"
	"log"
	cnf"simple_crude/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDB *mongo.Client
)

func Connect() error {
	cfg, err := cnf.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	clientOptions := options.Client().ApplyURI(cfg.Mongodb.Mongourl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("could not ping MongoDB: %v", err)
	}

	mongoDB = client
	log.Println("Connected to MongoDB successfully")
	return nil
}

func GetDB() *mongo.Client {
	return mongoDB
}

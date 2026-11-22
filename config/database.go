package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
Run mongo instance on docker
docker run -d --name mongodb -p 27017:27017 -v mongodb_data:/data/db mongo:8.0
*/

func ConnectDB() *mongo.Client {
	log.Println("Attempting to connect to MongoDB...")
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not reachable: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")
	return client
}

var Client = ConnectDB()

func OpenConnection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Please call ConnectDB first.")
	}
	return Client.Database("userdb").Collection(collectionName)
}

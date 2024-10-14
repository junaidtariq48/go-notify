package db

import (
	"context"
	"log"
	"notify/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongo initializes a MongoDB connection and returns the client
func InitMongo() *mongo.Client {
	// mongoURI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI).SetAuth(options.Credential{
		AuthSource: config.AppConfig.MongoDB,
		Username:   config.AppConfig.MongoUser,
		Password:   config.AppConfig.MongoPassword,
	})
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

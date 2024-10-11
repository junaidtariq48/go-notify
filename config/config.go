package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI   string
	RedisAddr  string
	RedisPass  string
	ServerPort string
}

// LoadConfig loads the config from environment variables
func LoadConfig() Config {
	return Config{
		MongoURI:   "mongodb://localhost:27017", // Change to env variable in real deployment
		RedisAddr:  "localhost:6379",            // Change to env variable
		RedisPass:  "",                          // Change to env variable
		ServerPort: ":8080",
	}
}

// InitMongo initializes a MongoDB connection and returns the client
func InitMongo(cfg Config) *mongo.Client {
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
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

// InitRedis initializes a Redis client
func InitRedis(cfg Config) *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // No password set
		DB:       0,             // Use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
	return rdb
}

func InitRabbitMQ() *amqp.Connection {
	rabbitmqURI := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	log.Println("Connected to RabbitMQ")
	return conn
}

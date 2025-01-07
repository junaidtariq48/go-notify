package main

import (
	"context"
	"log"
	"net/http"
	"notify/config"
	"notify/controllers"
	"notify/pkg/amqp"
	"notify/pkg/db"
	"notify/processors"
	"notify/workers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the logger and config
	config.InitLogger()

	// Load config
	config.InitializeConfig()

	// Initialize MongoDB, Redis, etc.
	db := db.InitMongo()
	defer db.Disconnect(context.Background())

	// redisClient := redis.InitRedis()

	// defer redisClient.Close()

	// Initialize RabbitMQ connection
	rabbitMQConn := amqp.InitRabbitMQ()
	defer rabbitMQConn.Close()

	rabbitMQChannel, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer rabbitMQChannel.Close()

	// Declare all the queues here
	amqp.DeclareQueue(rabbitMQChannel, workers.MainNotificationQueue)
	amqp.DeclareQueue(rabbitMQChannel, workers.EmailQueue)
	amqp.DeclareQueue(rabbitMQChannel, workers.SMSQueue)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create main notification processor
	mainProcessor := workers.NewMainNotificationWorker(rabbitMQChannel, db, config.Logger)
	// emailWorker := workers.NewNotificationWorker(rabbitMQChannel, db, workers.EmailQueue, processors.EmailProcessor)
	smsVerificationWorker := workers.NewNotificationWorker(rabbitMQChannel, db, workers.SMSQueue, processors.SMSProcessor)

	// Start workers
	go mainProcessor.Start(ctx)
	// go emailWorker.Start(ctx)
	go smsVerificationWorker.Start(ctx)

	// Initialize the router
	router := controllers.InitRouter(ctx, rabbitMQChannel, db)

	err1 := http.ListenAndServe(config.AppConfig.ServerPort, router)
	if err1 != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

	log.Printf("Server started on: %v", config.AppConfig.ServerPort)

	// Wait for shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	config.Logger.Info("Shutting down gracefully...")
	cancel() // This will stop all processors and workers
}

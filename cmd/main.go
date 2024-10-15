package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notify/config"
	"notify/controllers"
	"notify/pkg/db"
	"notify/pkg/redis"
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

	redisClient := redis.InitRedis()

	defer redisClient.Close()

	// Initialize RabbitMQ connection
	// rabbitMQConn := config.InitRabbitMQ()
	// defer rabbitMQConn.Close()

	// rabbitMQChannel, err := rabbitMQConn.Channel()
	// if err != nil {
	// 	log.Fatalf("Failed to open RabbitMQ channel: %s", err)
	// }
	// defer rabbitMQChannel.Close()

	// Start workers for each notification type
	// go workers.StartNotificationWorker(redisClient, db)

	// go workers.StartEmailWorker(redisClient, db)

	// 	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// emailWorker := NewNotificationWorker(redisClient, mongoClient, "email_queue", emailProcessor)
	// smsWorker := NewNotificationWorker(redisClient, mongoClient, "sms_queue", smsProcessor)

	// go emailWorker.Start(ctx)
	// go smsWorker.Start(ctx)

	// go workers.StartSMSWorker(redisClient, db)
	// go workers.StartPushWorker(redisClient, db)
	// go workers.StartWhatsAppWorker(redisClient, db)

	// Example payload as map
	// payload := map[string]interface{}{
	// 	"to":      "test@example.com",
	// 	"subject": "Welcome!",
	// 	"body":    "Welcome to our platform!",
	// }

	// // Convert map to JSON string
	// payloadJSON, errr := json.Marshal(payload)
	// if errr != nil {
	// 	log.Fatalf("Failed to marshal payload: %v", errr)
	// }

	// // Example Notification
	// emailNotification := models.Notification{
	// 	Type:     "email",
	// 	Provider: "sendgrid",
	// 	Payload:  string(payloadJSON),
	// }

	// logrus.Info(emailNotification)

	// // Enqueue Notification to Redis
	// err := queues.EnqueueNotification("email", emailNotification)
	// // err := queues.EnqueueNotification("email", emailNotification)
	// if err != nil {
	// 	log.Fatalf("Failed to enqueue notification: %v", err)
	// }

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create main notification processor
	mainProcessor := workers.NewMainNotificationWorker(redisClient, db, config.Logger)

	// Create workers
	// emailWorker := workers.NewNotificationWorker(redisClient, db, "email_queue", processors.EmailProcessor)
	emailWorker := workers.NewNotificationWorker(redisClient, db, workers.EmailQueue, processors.EmailProcessor)
	// smsWorker := workers.NewNotificationWorker(redisClient, db, "sms_queue", processors.SMSProcessor)

	// Start workers
	go mainProcessor.Start(ctx)
	go emailWorker.Start(ctx)
	// go smsWorker.Start(ctx)

	// Initialize HTTP server
	// router := controllers.InitRouter()
	// Initialize the router
	router := controllers.InitRouter(redisClient, db)

	err := http.ListenAndServe(config.AppConfig.ServerPort, router)
	fmt.Println(err)
	if err != nil {
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

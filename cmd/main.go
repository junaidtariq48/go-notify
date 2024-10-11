package main

import (
	"log"
	"net/http"
	"notify/config"
	"notify/controllers"
	"notify/workers"
)

func main() {
	// Initialize the logger
	config.InitLogger()
	// Load config
	cfg := config.LoadConfig()

	// Initialize MongoDB, Redis, etc.
	db := config.InitMongo(cfg)
	redisClient := config.InitRedis(cfg)

	// Initialize RabbitMQ connection
	// rabbitMQConn := config.InitRabbitMQ()
	// defer rabbitMQConn.Close()

	// rabbitMQChannel, err := rabbitMQConn.Channel()
	// if err != nil {
	// 	log.Fatalf("Failed to open RabbitMQ channel: %s", err)
	// }
	// defer rabbitMQChannel.Close()

	// Start workers for each notification type
	go workers.StartEmailWorker(redisClient, db)
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

	// Initialize HTTP server
	// router := controllers.InitRouter()
	// Initialize the router
	router := controllers.InitRouter(redisClient, db)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, router))
}

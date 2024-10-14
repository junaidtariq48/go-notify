package main

import (
	"fmt"
	"log"
	"net/http"
	"notify/config"
	"notify/controllers"
	"notify/pkg/db"
	"notify/pkg/redis"
	"notify/workers"
)

func main() {
	// Initialize the logger and config
	config.InitLogger()

	// Load config
	config.InitializeConfig()

	// Initialize MongoDB, Redis, etc.
	db := db.InitMongo()
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
	go workers.StartNotificationWorker(redisClient, db)

	go workers.StartEmailWorker(redisClient, db)

	// emailProcessor := func(ctx context.Context, notification models.Notification) error {
	// 	return services.SendEmail(ctx, notification)
	// }

	// smsProcessor := func(ctx context.Context, notification models.Notification) error {
	// 	return services.SendSMS(ctx, notification)
	// }

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

	err := http.ListenAndServe(config.AppConfig.ServerPort, router)
	fmt.Println(err)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

	log.Printf("Server started on: %v", config.AppConfig.ServerPort)
}

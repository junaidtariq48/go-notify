package controllers

import (
	"encoding/json"
	"net/http"
	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func CreateNotification(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, db *mongo.Client) {
	var notification models.Notification

	// Decode the incoming request payload into the Notification struct
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set initial status and timestamps
	notification.Status = "pending"
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	// Initialize the MongoDB repository
	repo := repositories.NewNotificationRepository(db)

	// Save the notification to MongoDB
	insertedID, err := repo.SaveNotification(&notification)
	if err != nil {
		config.Logger.WithError(err).Error("Failed to save notification to MongoDB")
		http.Error(w, "Failed to save notification", http.StatusInternalServerError)
		return
	}

	notification.ID = insertedID.Hex()
	config.Logger.WithField("notification_id", insertedID).Info("Notification saved to MongoDB")

	// Choose between Redis or RabbitMQ based on notification type (for example)
	var queueName string
	// var useRabbitMQ bool

	switch notification.Type {
	case "email":
		queueName = "email"
		// useRabbitMQ = false
	case "sms":
		queueName = "sms"
		// useRabbitMQ = true // For example, SMS can use RabbitMQ
	case "push":
		queueName = "push"
		// useRabbitMQ = false
	case "whatsapp":
		queueName = "whatsapp"
		// useRabbitMQ = true
	default:
		http.Error(w, "Invalid notification type", http.StatusBadRequest)
		return
	}

	// if useRabbitMQ {
	// 	// Enqueue to RabbitMQ
	// 	message, _ := json.Marshal(notification)
	// 	err = queues.EnqueueRabbitMQ(rabbitMQChannel, queueName, message)
	// } else {
	// Enqueue to Redis
	err = queues.EnqueueNotification(queueName, notification)
	// }

	if err != nil {
		config.Logger.WithError(err).Error("Failed to enqueue notification")
		http.Error(w, "Failed to enqueue notification", http.StatusInternalServerError)
		return
	}

	config.Logger.WithField("notification_id", insertedID).Info("Notification enqueued successfully")

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "success",
		"notification_id": insertedID,
	})
}

// func CreateNotification(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, db *mongo.Client, rabbitMQChannel *amqp.Channel) {
// 	var notification models.Notification

// 	// Decode the incoming request payload into the Notification struct
// 	err := json.NewDecoder(r.Body).Decode(&notification)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Set initial status and timestamps
// 	notification.Status = "pending"
// 	notification.CreatedAt = time.Now()
// 	notification.UpdatedAt = time.Now()

// 	// Initialize the MongoDB repository
// 	repo := repositories.NewNotificationRepository(db)

// 	// Save the notification to MongoDB
// 	insertedID, err := repo.SaveNotification(&notification)
// 	if err != nil {
// 		config.Logger.WithError(err).Error("Failed to save notification to MongoDB")
// 		http.Error(w, "Failed to save notification", http.StatusInternalServerError)
// 		return
// 	}

// 	config.Logger.WithField("notification_id", insertedID).Info("Notification saved to MongoDB")

// 	// Choose between Redis or RabbitMQ based on notification type (for example)
// 	var queueName string
// 	var useRabbitMQ bool

// 	switch notification.Type {
// 	case "email":
// 		queueName = "email"
// 		useRabbitMQ = false
// 	case "sms":
// 		queueName = "sms"
// 		useRabbitMQ = true // For example, SMS can use RabbitMQ
// 	case "push":
// 		queueName = "push"
// 		useRabbitMQ = false
// 	case "whatsapp":
// 		queueName = "whatsapp"
// 		useRabbitMQ = true
// 	default:
// 		http.Error(w, "Invalid notification type", http.StatusBadRequest)
// 		return
// 	}

// 	if useRabbitMQ {
// 		// Enqueue to RabbitMQ
// 		message, _ := json.Marshal(notification)
// 		err = queues.EnqueueRabbitMQ(rabbitMQChannel, queueName, message)
// 	} else {
// 		// Enqueue to Redis
// 		err = queues.EnqueueNotification(queueName, notification)
// 	}

// 	if err != nil {
// 		config.Logger.WithError(err).Error("Failed to enqueue notification")
// 		http.Error(w, "Failed to enqueue notification", http.StatusInternalServerError)
// 		return
// 	}

// 	config.Logger.WithField("notification_id", insertedID).Info("Notification enqueued successfully")

// 	// Return a success response
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"status":          "success",
// 		"notification_id": insertedID,
// 	})
// }

// func InitRouter() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/notifications", CreateNotification())
// 	return mux
// }

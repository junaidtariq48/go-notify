package services

import (
	"fmt"
	"notify/config"
	"notify/models"
	"notify/queues"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// SendEmail processes the email notification
func ProcessNotification(db *mongo.Client, redisClient *redis.Client, notification models.Notification) error {
	// Choose between Redis or RabbitMQ based on notification type (for example)
	var queueName string

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
		return fmt.Errorf("unsupported notificaiton type: %s", notification.Type)
	}

	err := queues.EnqueueNotification(queueName, notification)

	if err != nil {
		config.Logger.WithError(err).Error("Failed to enqueue notification")
		return fmt.Errorf("Failed to enqueue notification: %s", notification.Type)
	}

	config.Logger.WithField("notification_id", notification.ID).Info("Notification enqueued successfully")

	return nil
}

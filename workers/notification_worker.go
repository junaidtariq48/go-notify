package workers

import (
	"notify/config"
	"notify/queues"
	"notify/repositories"
	"notify/services"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// StartNotificationWorker processes email notifications from the Redis queue
func StartNotificationWorker(redisClient *redis.Client, db *mongo.Client) {
	repo := repositories.NewNotificationRepository(db)
	for {
		notification, err := queues.DequeueNotification("notifications_queue")
		if err != nil {
			config.Logger.WithError(err).Error("Error dequeuing notifications")
			time.Sleep(1 * time.Second)
			continue
		}

		if notification == nil {
			// Queue is empty, no need to process, just sleep and retry
			time.Sleep(1 * time.Second)
			continue
		}

		// Set initial status and timestamps
		notification.Status = "pending"
		notification.CreatedAt = time.Now()
		notification.UpdatedAt = time.Now()

		// Save the notification to MongoDB
		insertedID, err := repo.SaveNotification(notification)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"type":         notification.Type,
				"notification": notification,
			}).Error("Error Processing notification")
		}

		notification.ID = insertedID.Hex()

		config.Logger.WithFields(logrus.Fields{
			"ID":       notification.ID,
			"type":     notification.Type,
			"provider": notification.Provider,
		}).Info("Processing new notification")

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
			config.Logger.WithError(err).Error("unsupported notificaiton type")
		}

		err = queues.EnqueueNotification(queueName, *notification)

		if err != nil {
			config.Logger.WithError(err).Error("Failed to enqueue notification")
		}

		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"queue_name":      queueName,
		}).Info("Notification enqueued successfully")

		// process the notification
		// err = services.ProcessNotification(db, redisClient, *notification)
		// if err != nil {
		// 	config.Logger.WithFields(logrus.Fields{
		// 		"notification_id": notification.ID,
		// 		"error":           err,
		// 	}).Error("Failed to process notification")
		// 	// _ = repo.UpdateNotificationStatus(notification.ID, "failed")
		// } else {
		// 	config.Logger.WithFields(logrus.Fields{
		// 		"notification_id": notification.ID,
		// 	}).Info("Notification process successfully", notification.ID)
		// 	// _ = repo.UpdateNotificationStatus(notification.ID, "success")
		// }
	}
}

// StartEmailWorker processes email notifications from the Redis queue
func StartEmailWorker(redisClient *redis.Client, db *mongo.Client) {
	repo := repositories.NewNotificationRepository(db)

	for {
		notification, err := queues.DequeueNotification("email")
		if err != nil {
			config.Logger.WithError(err).Error("Error dequeuing email notification")
			time.Sleep(1 * time.Second)
			continue
		}

		if notification == nil {
			// Queue is empty, no need to process, just sleep and retry
			time.Sleep(1 * time.Second)
			continue
		}

		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"provider":        notification.Provider,
		}).Info("Processing email notification")

		// Send email using the entire notification object
		err = services.SendEmail(*notification)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
				"error":           err,
			}).Error("Failed to send email")
			_ = repo.UpdateNotificationStatus(notification.ID, "failed")
		} else {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
			}).Info("Email sent successfully", notification.ID)
			_ = repo.UpdateNotificationStatus(notification.ID, "success")
		}
	}
}

func StartSMSWorker(redisClient *redis.Client, db *mongo.Client) {
	repo := repositories.NewNotificationRepository(db)

	for {
		notification, err := queues.DequeueNotification("sms")
		if err != nil {
			config.Logger.WithError(err).Error("Error dequeuing SMS notification")
			time.Sleep(1 * time.Second)
			continue
		}

		if notification == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"provider":        notification.Provider,
		}).Info("Processing SMS notification")

		// Send SMS using the entire notification object
		err = services.SendSMS(*notification)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
				"error":           err,
			}).Error("Failed to send SMS")
			_ = repo.UpdateNotificationStatus(notification.ID, "failed")
		} else {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
			}).Info("SMS sent successfully")
			_ = repo.UpdateNotificationStatus(notification.ID, "success")
		}
	}
}

func StartPushWorker(redisClient *redis.Client, db *mongo.Client) {
	repo := repositories.NewNotificationRepository(db)

	for {
		notification, err := queues.DequeueNotification("push")
		if err != nil {
			config.Logger.WithError(err).Error("Error dequeuing push notification")
			time.Sleep(1 * time.Second)
			continue
		}

		if notification == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"provider":        notification.Provider,
		}).Info("Processing push notification")

		err = services.SendPushNotification(*notification)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
				"error":           err,
			}).Error("Failed to send push notification")
			_ = repo.UpdateNotificationStatus(notification.ID, "failed")
		} else {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
			}).Info("Push notification sent successfully")
			_ = repo.UpdateNotificationStatus(notification.ID, "success")
		}
	}
}

func StartWhatsAppWorker(redisClient *redis.Client, db *mongo.Client) {
	repo := repositories.NewNotificationRepository(db)

	for {
		notification, err := queues.DequeueNotification("whatsapp")
		if err != nil {
			config.Logger.WithError(err).Error("Error dequeuing WhatsApp notification")
			time.Sleep(1 * time.Second)
			continue
		}

		if notification == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"provider":        notification.Provider,
		}).Info("Processing WhatsApp notification")

		err = services.SendWhatsApp(*notification)
		if err != nil {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
				"error":           err,
			}).Error("Failed to send WhatsApp message")
			_ = repo.UpdateNotificationStatus(notification.ID, "failed")
		} else {
			config.Logger.WithFields(logrus.Fields{
				"notification_id": notification.ID,
			}).Info("WhatsApp message sent successfully")
			_ = repo.UpdateNotificationStatus(notification.ID, "success")
		}
	}
}

// // Worker to process Redis and RabbitMQ queues
// func StartWorker(redisClient *redis.Client, rabbitMQChannel *amqp.Channel, queueName string, notificationType string, useRabbitMQ bool) {
// 	if useRabbitMQ {
// 		// RabbitMQ queue worker
// 		log.Printf("Worker started for RabbitMQ queue: %s", queueName)
// 		msgs, err := queues.DequeueRabbitMQ(rabbitMQChannel, queueName)
// 		if err != nil {
// 			log.Fatalf("Failed to consume from RabbitMQ queue: %s", err)
// 		}

// 		for msg := range msgs {
// 			var notification models.Notification
// 			err := json.Unmarshal(msg.Body, &notification)
// 			if err != nil {
// 				log.Printf("Failed to unmarshal message: %s", err)
// 				continue
// 			}

// 			processNotification(notification, notificationType)
// 		}
// 	} else {
// 		// Redis queue worker
// 		log.Printf("Worker started for Redis queue: %s", queueName)
// 		for {
// 			notification, err := queues.DequeueNotification(redisClient, queueName)
// 			if err != nil {
// 				log.Printf("Error dequeuing notification: %s", err)
// 				continue
// 			}

// 			processNotification(notification, notificationType)
// 		}
// 	}
// }

// // Function to process the notification based on type
// func processNotification(notification *models.Notification, notificationType string) {
// 	switch notificationType {
// 	case "email":
// 		err := services.SendEmail(*notification)
// 		if err != nil {
// 			log.Printf("Failed to send email: %s", err)
// 		} else {
// 			log.Println("Email sent successfully")
// 		}
// 	case "sms":
// 		err := services.SendSMS(*notification)
// 		if err != nil {
// 			log.Printf("Failed to send SMS: %s", err)
// 		} else {
// 			log.Println("SMS sent successfully")
// 		}
// 	case "push":
// 		err := services.SendPushNotification(*notification)
// 		if err != nil {
// 			log.Printf("Failed to send push notification: %s", err)
// 		} else {
// 			log.Println("Push notification sent successfully")
// 		}
// 	case "whatsapp":
// 		err := services.SendWhatsApp(*notification)
// 		if err != nil {
// 			log.Printf("Failed to send WhatsApp message: %s", err)
// 		} else {
// 			log.Println("WhatsApp message sent successfully")
// 		}
// 	default:
// 		log.Printf("Unknown notification type: %s", notificationType)
// 	}
// }

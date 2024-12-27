package workers

import (
	"context"
	"encoding/json"
	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationProcessor func(context.Context, repositories.Repositories, models.Notification) error

type NotificationWorker struct {
	channel   *amqp.Channel
	db        *mongo.Client
	queueName string
	processor NotificationProcessor
	repo      *repositories.Repositories
}

func NewNotificationWorker(channel *amqp.Channel, db *mongo.Client, queueName string, processor NotificationProcessor) *NotificationWorker {
	return &NotificationWorker{
		channel:   channel,
		db:        db,
		queueName: queueName,
		processor: processor,
		repo:      repositories.NewRepositories(db),
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {

	config.Logger.Info("Started %s queue..", w.queueName)

	messages, err := queues.ConsumeNotification(ctx, w.channel, w.queueName)
	if err != nil {
		config.Logger.WithError(err).Fatal("Failed to start consuming notifications")
	}

	for {
		select {
		case <-ctx.Done():
			config.Logger.Info("Stopping %s notification worker", w.queueName)
			return
		case msg := <-messages:
			var notification models.Notification
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				config.Logger.WithError(err).Error("Failed to parse notification")
				continue
			}
			w.processNotification(ctx, &notification)
		}
	}
}

// func (w *NotificationWorker) Start(ctx context.Context) {
// 	config.Logger.Info("Starting email notification processor")

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			w.processNotification(ctx)
// 		}
// 	}
// }

func (w *NotificationWorker) processNotification(ctx context.Context, notification *models.Notification) {
	// notification, err := queues.DequeueNotification(ctx, w.redisClient, w.queueName)
	// if err != nil {
	// 	if err != redis.Nil {
	// 		config.Logger.WithError(err).Errorf("Error dequeuing %s notification", w.queueName)
	// 	}
	// 	time.Sleep(1 * time.Second)
	// 	return
	// }

	// if notification == nil {
	// 	time.Sleep(1 * time.Second)
	// 	return
	// }

	// config.Logger.WithFields(logrus.Fields{
	// 	"notification_id": notification.ID,
	// 	"provider":        notification.Provider,
	// }).Infof("Processing %s notification", w.queueName)

	// err = w.processor(ctx, *w.repo, *notification)
	// if err != nil {
	// 	config.Logger.WithFields(logrus.Fields{
	// 		"notification_id": notification.ID,
	// 		"error":           err,
	// 	}).Errorf("Failed to process %s notification", w.queueName)
	// 	_ = w.repo.NotificationRepo.UpdateNotificationStatus(ctx, notification.ID, "failed")
	// } else {
	// 	config.Logger.WithFields(logrus.Fields{
	// 		"notification_id": notification.ID,
	// 	}).Infof("%s notification processed successfully", w.queueName)
	// 	_ = w.repo.NotificationRepo.UpdateNotificationStatus(ctx, notification.ID, "success")
	// }
}

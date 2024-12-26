package workersRabbit

import (
	"context"
	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"
	"time"

	"github.com/sirupsen/logrus"
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
	config.Logger.Info("Starting notification processor")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.processNotification(ctx)
		}
	}
}

func (w *NotificationWorker) processNotification(ctx context.Context) {
	// Dequeue notification from RabbitMQ
	notification, err := queues.DequeueRabbitNotification(ctx, w.channel, w.queueName)
	if err != nil {
		config.Logger.WithError(err).Errorf("Error dequeuing %s notification", w.queueName)
		time.Sleep(10 * time.Second)
		return
	}

	if notification == nil {
		time.Sleep(1 * time.Second)
		return
	}

	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"provider":        notification.Provider,
	}).Infof("Processing %s notification", w.queueName)

	// Process the notification
	err = w.processor(ctx, *w.repo, *notification)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
			"error":           err,
		}).Errorf("Failed to process %s notification", w.queueName)
		_ = w.repo.NotificationRepo.UpdateNotificationStatus(ctx, notification.ID, "failed")
	} else {
		config.Logger.WithFields(logrus.Fields{
			"notification_id": notification.ID,
		}).Infof("%s notification processed successfully", w.queueName)
		_ = w.repo.NotificationRepo.UpdateNotificationStatus(ctx, notification.ID, "success")
	}
}

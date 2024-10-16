package workers

import (
	"context"
	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationProcessor func(context.Context, models.Notification) error

type NotificationWorker struct {
	redisClient *redis.Client
	db          *mongo.Client
	queueName   string
	processor   NotificationProcessor
	repo        *repositories.Repositories
}

func NewNotificationWorker(redisClient *redis.Client, db *mongo.Client, queueName string, processor NotificationProcessor) *NotificationWorker {
	return &NotificationWorker{
		redisClient: redisClient,
		db:          db,
		queueName:   queueName,
		processor:   processor,
		repo:        repositories.NewRepositories(db),
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	config.Logger.Info("Starting email notification processor")

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
	notification, err := queues.DequeueNotification(ctx, w.redisClient, w.queueName)
	if err != nil {
		if err != redis.Nil {
			config.Logger.WithError(err).Errorf("Error dequeuing %s notification", w.queueName)
		}
		time.Sleep(1 * time.Second)
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

	err = w.processor(ctx, *notification)
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

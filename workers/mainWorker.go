package workers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"

	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"
)

const (
	MainNotificationQueue = "notifications_queue"
	EmailQueue            = "email_queue"
	SMSQueue              = "sms_queue"
)

type MainNotificationWorker struct {
	channel *amqp.Channel
	logger  *logrus.Logger
	db      *mongo.Client
	repo    *repositories.Repositories
}

func NewMainNotificationWorker(channel *amqp.Channel, db *mongo.Client, logger *logrus.Logger) *MainNotificationWorker {
	return &MainNotificationWorker{
		channel: channel,
		logger:  logger,
		db:      db,
		repo:    repositories.NewRepositories(db),
	}
}

func (p *MainNotificationWorker) Start(ctx context.Context) {
	p.logger.Info("Starting main notification worker")

	messages, err := queues.ConsumeNotification(ctx, p.channel, MainNotificationQueue)
	if err != nil {
		p.logger.WithError(err).Fatal("Failed to start consuming notifications")
	}

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Stopping main notification worker")
			return
		case msg := <-messages:
			var notification models.Notification
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				p.logger.WithError(err).Error("Failed to parse notification")
				continue
			}
			p.processNotification(ctx, &notification)
		}
	}
}

func (p *MainNotificationWorker) processNotification(ctx context.Context, notification *models.Notification) {
	notification.Status = "pending"
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	insertedID, err := p.repo.NotificationRepo.SaveNotification(ctx, notification)
	if err != nil {
		config.Logger.WithError(err).Error("Error saving notification")
		return
	}

	notification.ID = insertedID.Hex()
	p.distributeNotification(ctx, notification)
}

func (p *MainNotificationWorker) distributeNotification(ctx context.Context, notification *models.Notification) {
	var targetQueue string
	switch notification.Type {
	case "email":
		targetQueue = EmailQueue
	case "sms":
		targetQueue = SMSQueue
	default:
		p.logger.WithField("type", notification.Type).Error("Unknown notification type")
		return
	}

	err := queues.PublishNotification(ctx, p.channel, targetQueue, *notification)
	if err != nil {
		p.logger.WithError(err).WithField("queue", targetQueue).Error("Failed to publish notification to queue")
	}
}

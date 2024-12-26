package workersRabbit

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/repositories"

	"github.com/streadway/amqp"
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
	p.logger.Info("Starting main notification Worker listening to rabbit....")

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Stopping main notification Worker")
			return
		default:
			p.processMainQueue(ctx)
		}
	}
}

func (p *MainNotificationWorker) processMainQueue(ctx context.Context) {

	notification, err := queues.DequeueRabbitNotification(ctx, p.channel, MainNotificationQueue)
	if err != nil {
		p.logger.WithError(err).Error("Error consuming from main notification queue")
		return
	}

	if notification == nil {
		// fmt.Println("No notification in the queue yet")
		return
	}

	// Set initial status and timestamps
	notification.Status = "pending"
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	// Save the notification to MongoDB
	insertedID, err := p.repo.NotificationRepo.SaveNotification(ctx, notification)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"type":         notification.Type,
			"notification": notification,
		}).Error("Error Processing notification")
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

	err := queues.EnqueueRabbitNotification(ctx, p.channel, targetQueue, *notification)
	if err != nil {
		p.logger.WithError(err).WithField("queue", targetQueue).Error("Error pushing to specific queue")
		return
	}

	p.logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"type":            notification.Type,
		"queue":           targetQueue,
	}).Info("Notification distributed to specific queue")
}

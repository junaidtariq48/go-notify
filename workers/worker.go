package workers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"notify/config"
	"notify/models"
	"notify/repositories"
)

const (
	MainNotificationQueue = "queue:notifications_queue"
	EmailQueue            = "email_queue"
	SMSQueue              = "sms_queue"
)

type MainNotificationWorker struct {
	redisClient *redis.Client
	logger      *logrus.Logger
	db          *mongo.Client
	repo        *repositories.Repositories
}

func NewMainNotificationWorker(redisClient *redis.Client, db *mongo.Client, logger *logrus.Logger) *MainNotificationWorker {
	return &MainNotificationWorker{
		redisClient: redisClient,
		logger:      logger,
		db:          db,
		repo:        repositories.NewRepositories(db),
	}
}

func (p *MainNotificationWorker) Start(ctx context.Context) {
	p.logger.Info("Starting main notification Worker")

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

	result, err := p.redisClient.BLPop(ctx, 0, MainNotificationQueue).Result()
	if err != nil {
		if err != redis.Nil {
			p.logger.WithError(err).Error("Error dequeuing from main notification queue")
		}
		return
	}

	if len(result) != 2 {
		p.logger.Error("Unexpected result format from Redis")
		return
	}

	var notification models.Notification
	err = json.Unmarshal([]byte(result[1]), &notification)
	if err != nil {
		p.logger.WithError(err).Error("Error unmarshaling notification")
		return
	}

	// Set initial status and timestamps
	notification.Status = "pending"
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	// Save the notification to MongoDB
	insertedID, err := p.repo.NotificationRepo.SaveNotification(ctx, &notification)
	if err != nil {
		config.Logger.WithFields(logrus.Fields{
			"type":         notification.Type,
			"notification": notification,
		}).Error("Error Processing notification")
	}

	notification.ID = insertedID.Hex()

	p.distributeNotification(ctx, &notification)
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

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		p.logger.WithError(err).Error("Error marshaling notification for distribution")
		return
	}

	err = p.redisClient.RPush(ctx, targetQueue, notificationJSON).Err()
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

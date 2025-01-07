package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	EmailRepo        *EmailRepository
	SmsRepo          *SmsRepository
	NotificationRepo *NotificationRepository
}

// NewRepositories creates a new instance of Repositories
func NewRepositories(db *mongo.Client) *Repositories {
	return &Repositories{
		EmailRepo:        NewEmailRepository(db),
		SmsRepo:          NewSmsRepository(db),
		NotificationRepo: NewNotificationRepository(db),
	}
}

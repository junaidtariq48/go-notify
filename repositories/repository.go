package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	EmailRepo        *EmailRepository
	NotificationRepo *NotificationRepository
}

// NewRepositories creates a new instance of Repositories
func NewRepositories(db *mongo.Client) *Repositories {
	return &Repositories{
		EmailRepo:        NewEmailRepository(db),
		NotificationRepo: NewNotificationRepository(db),
	}
}

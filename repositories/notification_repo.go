package repositories

import (
	"context"
	"log"
	"notify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	Collection *mongo.Collection
}

// NewNotificationRepository creates a new instance of NotificationRepository
func NewNotificationRepository(db *mongo.Client) *NotificationRepository {
	collection := db.Database("notifications_db").Collection("notifications")
	return &NotificationRepository{Collection: collection}
}

// SaveNotification saves a new notification in the MongoDB
func (r *NotificationRepository) SaveNotification(notification *models.Notification) (primitive.ObjectID, error) {
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	result, err := r.Collection.InsertOne(context.TODO(), notification)
	if err != nil {
		log.Printf("Error inserting notification into MongoDB: %v", err)
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateNotificationStatus updates the status of a notification
func (r *NotificationRepository) UpdateNotificationStatus(id string, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Failed to update notification status:", err)
	}
	return err
}

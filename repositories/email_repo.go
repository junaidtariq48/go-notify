package repositories

import (
	"context"
	"fmt"
	"log"
	"notify/config"
	constants "notify/contants"
	"notify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailRepository struct {
	Collection *mongo.Collection
}

// NewNotificationRepository creates a new instance of NotificationRepository
func NewEmailRepository(db *mongo.Client) *EmailRepository {
	collection := db.Database(config.AppConfig.MongoDB).Collection(constants.EMAIL_COLLECTION)
	return &EmailRepository{Collection: collection}
}

// SaveEmail saves a new notification in the MongoDB
func (r *EmailRepository) SaveEmail(ctx context.Context, email *models.Email) (primitive.ObjectID, error) {
	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	result, err := r.Collection.InsertOne(context.TODO(), email)
	if err != nil {
		log.Printf("Error inserting email into MongoDB: %v", err)
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateEmailStatus updates the status of a notification
func (r *EmailRepository) UpdateEmailStatus(id string, status string) error {
	objectID, errr := primitive.ObjectIDFromHex(id)
	if errr != nil {
		log.Fatal(errr)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	fmt.Println(filter)
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Failed to update email status:", err)
	}
	return err
}

func (r *EmailRepository) UpdateEmailResposne(id string, response string) error {
	objectID, errr := primitive.ObjectIDFromHex(id)
	if errr != nil {
		log.Fatal(errr)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"response":   response,
			"updated_at": time.Now(),
		},
	}
	fmt.Println(filter)
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Failed to update email response:", err)
	}
	return err
}

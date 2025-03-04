package repositories

import (
	"context"
	"log"
	"notify/config"
	constants "notify/contants"
	"notify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SmsRepository struct {
	Collection *mongo.Collection
}

// NewNotificationRepository creates a new instance of NotificationRepository
func NewSmsRepository(db *mongo.Client) *SmsRepository {
	collection := db.Database(config.AppConfig.MongoDB).Collection(constants.SMS_COLLECTION)
	return &SmsRepository{Collection: collection}
}

// SaveEmail saves a new notification in the MongoDB
func (r *SmsRepository) SaveSMS(ctx context.Context, sms *models.SMSNotification) (primitive.ObjectID, error) {
	result, err := r.Collection.InsertOne(context.TODO(), sms)
	if err != nil {
		log.Printf("Error inserting sms verification into MongoDB: %v", err)
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateEmailStatus updates the status of a notification
func (r *SmsRepository) UpdateSMSStatus(ctx context.Context, id string, status string) error {
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

	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Failed to update email status:", err)
	}
	return err
}

func (r *SmsRepository) UpdateSMSResponse(ctx context.Context, id string, response string, status string) error {
	objectID, errr := primitive.ObjectIDFromHex(id)
	if errr != nil {
		log.Fatal(errr)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"response":   response,
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Failed to update sms response:", err)
	}
	return err
}

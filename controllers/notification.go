package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"notify/config"
	"notify/models"
	"notify/queues"
	"notify/workers"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func CreateNotification(w http.ResponseWriter, r *http.Request, ctx context.Context, channel *amqp.Channel, db *mongo.Client) {
	var notification models.Notification

	// Decode the incoming request payload into the Notification struct
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err1 := queues.PublishNotification(ctx, channel, workers.MainNotificationQueue, notification)
	if err1 != nil {
		config.Logger.WithError(err).WithField("queue", workers.MainNotificationQueue).Error("Failed to publish notification to queue")
	}

	config.Logger.WithField("Notification Type:", notification.Type).Info("Notification enqueued successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}

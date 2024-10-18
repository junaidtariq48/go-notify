package controllers

import (
	"encoding/json"
	"net/http"
	"notify/config"
	"notify/models"
	"notify/queues"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func CreateNotification(w http.ResponseWriter, r *http.Request, redisClient *redis.Client, db *mongo.Client) {
	var notification models.Notification

	// Decode the incoming request payload into the Notification struct
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = queues.EnqueueNotification(r.Context(), redisClient, "notifications_queue", notification)

	if err != nil {
		config.Logger.WithError(err).Error("Failed to enqueue notification")
		http.Error(w, "Failed to enqueue notification", http.StatusInternalServerError)
		return
	}

	config.Logger.WithField("notification_id", notification.Type).Info("Notification enqueued successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
	})
}

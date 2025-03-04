package models

import "time"

type SMSNotification struct {
	ID             string    `json:"id" bson:"_id,omitempty"` // MongoDB ID
	NotificationID string    `json:"notification_id" bson:"notification_id"`
	Recipient      string    `json:"recipient" bson:"recipient"` // Phone number
	Message        string    `json:"message" bson:"message"`     // SMS content
	Provider       string    `json:"provider" bson:"provider"`   // E.g., "Twilio", "Nexmo"
	Status         string    `json:"status" bson:"status"`       // "pending", "sent", "failed", etc.
	Type           string    `json:"type" bson:"type"`           //"code","general"
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

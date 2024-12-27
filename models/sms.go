package models

type SMSNotification struct {
	NotificationID string `json:"notification_id" bson:"notification_id"`
	Recipient      string `json:"recipient" bson:"recipient"` // Phone number
	Message        string `json:"message" bson:"message"`     // SMS content
	Provider       string `json:"provider" bson:"provider"`   // E.g., "Twilio", "Nexmo"
}

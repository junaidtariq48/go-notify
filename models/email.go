package models

import "time"

type EmailNotification struct {
	ID             string            `json:"id" bson:"_id,omitempty"`
	NotificationID string            `json:"notification_id" bson:"notification_id"`
	Recipient      string            `json:"recipient" bson:"recipient"`           // Email address
	RecipientName  string            `json:"recipient_name" bson:"recipient_name"` // customer name
	Subject        string            `json:"subject" bson:"subject"`               // Email subject
	TemplateID     string            `json:"template_id" bson:"template_id"`       // SendGrid dynamic template ID
	DynamicData    map[string]string `json:"dynamic_data" bson:"dynamic_data"`     // Variables for the template
	Status         string            `json:"status" bson:"status"`                 // "pending", "sent", "failed", etc.
	Type           string            `json:"type" bson:"type"`                     //"code","general"
	Provider       string            `json:"provider" bson:"provider"`             // E.g., "Twilio", "Nexmo"
	CreatedAt      time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" bson:"updated_at"`
}

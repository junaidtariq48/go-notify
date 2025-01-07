package models

import (
	"time"
)

type Notification struct {
	ID            string                 `json:"id" bson:"_id,omitempty"`    // MongoDB ID
	Type          string                 `json:"type" bson:"type"`           // "email" or "sms"
	Provider      string                 `json:"provider" bson:"provider"`   // "twilio" or "nexmo"
	Status        string                 `json:"status" bson:"status"`       // "pending", "sent", "failed", etc.
	Recipient     string                 `json:"recipient" bson:"recipient"` // Email address or phone number
	Payload       map[string]interface{} `json:"payload" bson:"payload"`     // Custom data (for email templates, SMS, etc.)
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" bson:"updated_at"`
	SentAt        *time.Time             `json:"sent_at,omitempty" bson:"sent_at,omitempty"`               // Optional
	FailureReason string                 `json:"failure_reason,omitempty" bson:"failure_reason,omitempty"` // If failed
}

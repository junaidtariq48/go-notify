package models

import (
	"time"
)

type Notification struct {
	ID        string    `bson:"_id,omitempty"`
	Type      string    `bson:"type"`
	Provider  string    `bson:"provider"`
	Payload   string    `bson:"payload"`
	Status    string    `bson:"status"` // e.g., "pending", "success", "failed"
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

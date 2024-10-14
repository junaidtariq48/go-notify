package models

import (
	"time"
)

type Email struct {
	ID              string    `bson:"_id,omitempty"`
	Notificaiton_id string    `bson:"notification_id"`
	To              string    `bson:"to"`
	From            string    `bson:"from"`
	Subject         string    `bson:"Subject"`
	Body            string    `bson:"Body"`
	response        string    `bson:"Response"`
	Status          string    `bson:"status"`
	CreatedAt       time.Time `bson:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

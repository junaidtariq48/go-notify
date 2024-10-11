package push_providers

import (
	"encoding/json"
	"log"
	"notify/models"
)

type FirebaeProvider struct{}

// Firebase processes the push notification using the Firebase provider
func (f *FirebaeProvider) Send(notification models.Notification) error {
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(notification.Payload), &payload)
	if err != nil {
		return err
	}

	log.Printf("Firebase: Sending push notification to %s with message %s", payload["user"], payload["message"])

	// Simulate sending push notification via Firebase
	// Normally here you'd call Firebase's API to send the push notification
	return nil
}

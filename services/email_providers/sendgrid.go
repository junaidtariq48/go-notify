package email_providers

import (
	"encoding/json"
	"log"
	"notify/models"
)

type SendGridProvider struct{}

// SendGrid processes the email notification using the SendGrid provider
func (s *SendGridProvider) Send(notification models.Notification) error {
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(notification.Payload), &payload)
	if err != nil {
		return err
	}

	log.Printf("SendGrid: Sending email to %s with subject %s", payload["to"], payload["subject"])

	// Simulate sending email via SendGrid
	// Normally here you'd call SendGrid's API to send the email
	return nil
}

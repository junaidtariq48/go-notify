package sms_providers

import (
	"encoding/json"
	"log"
	"notify/models"
)

// TwilioProvider sends SMS using Twilio
type TwilioProvider struct{}

// Twilio processes the SMS notification using the Twilio provider
func (t *TwilioProvider) Send(notification models.Notification) error {
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(notification.Payload), &payload)
	if err != nil {
		return err
	}

	log.Printf("Twilio: Sending SMS to %s with message %s", payload["to"], payload["message"])

	// Simulate sending SMS via Twilio
	// Normally here you'd call Twilio's API to send the SMS
	return nil
}

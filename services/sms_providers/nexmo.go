package sms_providers

import (
	"log"
	"notify/models"
)

// NexmoProvider sends SMS using Twilio
type NexmoProvider struct{}

// Nexmo processes the SMS notification using the Nexmo provider
func (n *NexmoProvider) Send(notification models.SMSNotification) error {
	var payload map[string]interface{}
	// err := json.Unmarshal([]byte(notification.Payload), &payload)
	// if err != nil {
	// 	return err
	// }

	log.Printf("Twilio: Sending SMS to %s with message %s", payload["to"], payload["message"])

	// Simulate sending SMS via Twilio
	// Normally here you'd call Twilio's API to send the SMS
	return nil
}

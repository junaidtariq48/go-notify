package whatsapp_providers

import (
	"log"
	"notify/models"
)

// TwilioWhatsAppProvider sends WhatsApp messages using Twilio
type TwilioWhatsAppProvider struct{}

// Twilio processes the WhatsApp notification using the Twilio provider
func (t *TwilioWhatsAppProvider) Send(notification models.Notification) error {
	var payload map[string]interface{}
	// err := json.Unmarshal([]byte(notification.Payload), &payload)
	// if err != nil {
	// 	return err
	// }

	log.Printf("Twilio: Sending WhatsApp message to %s with content %s", payload["to"], payload["content"])

	// Simulate sending WhatsApp message via Twilio
	// Normally here you'd call Twilio's API to send the WhatsApp message
	return nil
}

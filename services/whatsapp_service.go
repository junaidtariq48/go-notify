package services

import (
	"notify/config"
	"notify/models"

	"github.com/sirupsen/logrus"
)

// WhatsAppProvider defines the interface for WhatsApp providers
type WhatsAppProvider interface {
	Send(notification models.Notification) error
}

// SendWhatsApp sends a WhatsApp message using the appropriate provider
func SendWhatsApp(notification models.Notification) error {
	// provider := &TwilioWhatsAppProvider{} // We are using Twilio for WhatsApp messages
	var provider WhatsAppProvider

	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		// "provider":        notification.Provider,
	}).Info("Sending WhatsApp notification")

	// Select the WhatsApp provider based on notification.Provider
	// switch notification.Provider {
	// case "twilio":
	// 	provider = &whatsapp_providers.TwilioWhatsAppProvider{}
	// // case "another_whatsapp_provider":
	// //     return whatsapp_providers.AnotherWhatsAppProvider(notification)
	// default:
	// 	log.Println("Unsupported WhatsApp provider:", notification.Provider)
	// 	return fmt.Errorf("unsupported WhatsApp provider: %s", notification.Provider)
	// }

	return provider.Send(notification)
}

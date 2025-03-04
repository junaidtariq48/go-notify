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
		"provider":        notification.Provider,
	}).Info("Sending WhatsApp notification")

	return provider.Send(notification)
}

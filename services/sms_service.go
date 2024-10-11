package services

import (
	"fmt"
	"log"
	"notify/config"
	"notify/models"
	"notify/services/sms_providers"

	"github.com/sirupsen/logrus"
)

// SMSProvider defines the interface for SMS providers
type SMSProvider interface {
	Send(notification models.Notification) error
}

// SendSMS sends an SMS using the appropriate provider based on the country or other logic
func SendSMS(notification models.Notification) error {
	var provider SMSProvider

	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"provider":        notification.Provider,
	}).Info("Sending SMS notification")

	// Select the SMS provider based on notification.Provider
	switch notification.Provider {
	case "twilio":
		provider = &sms_providers.TwilioProvider{}
	case "nexmo":
		provider = &sms_providers.NexmoProvider{}
	default:
		log.Println("Unsupported SMS provider:", notification.Provider)
		return fmt.Errorf("unsupported SMS provider: %s", notification.Provider)
	}

	return provider.Send(notification)
}

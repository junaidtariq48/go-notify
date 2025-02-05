package services

import (
	"context"
	"fmt"
	"log"
	"notify/config"
	"notify/models"
	"notify/services/sms_providers"

	"github.com/sirupsen/logrus"
)

// SMSProvider defines the interface for SMS providers
type SMSProvider interface {
	Send(notification models.SMSNotification) ([]byte, error)
}

// SendSMS sends an SMS using the appropriate provider based on the country or other logic
func SendSMS(ctx context.Context, notification models.SMSNotification) ([]byte, error) {
	var provider SMSProvider

	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
	}).Info("Sending SMS notification")

	// Select the SMS provider based on notification.Provider
	switch notification.Provider {
	case "twilio":
		provider = &sms_providers.TwilioProvider{}
	case "nexmo":
		// provider = &sms_providers.NexmoProvider{}
	default:
		// provider = &sms_providers.TwilioProvider{} // default provider is twilio
		log.Println("Unsupported SMS provider:", notification.Provider)
		return nil, fmt.Errorf("unsupported SMS provider: %s", notification.Provider)
	}

	return provider.Send(notification)
}

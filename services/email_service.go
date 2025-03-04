package services

import (
	"context"
	"fmt"
	"log"
	"notify/config"
	"notify/models"
	"notify/services/email_providers"

	"github.com/sirupsen/logrus"
)

type EmailProvider interface {
	Send(ctx context.Context, notification models.EmailNotification) ([]byte, error)
}

// SendEmail processes the email notification
func SendEmail(ctx context.Context, notification models.EmailNotification) ([]byte, error) {
	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"provider":        notification.Provider,
	}).Info("Sending email notification")

	var provider EmailProvider

	// Select the email provider based on notification.Provider
	switch notification.Provider {
	case "sendgrid":
		provider = &email_providers.SendGridProvider{}
	default:
		log.Println("Unsupported email provider:", notification.Provider)
		return nil, fmt.Errorf("unsupported email provider: %s", notification.Provider)
	}

	return provider.Send(ctx, notification)
}

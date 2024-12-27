package services

import (
	"context"
	"notify/config"
	"notify/models"

	"github.com/sirupsen/logrus"
)

type EmailProvider interface {
	Send(ctx context.Context, notification models.Notification) error
}

// SendEmail processes the email notification
func SendEmail(ctx context.Context, notification models.Notification) error {
	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		// "provider":        notification.Provider,
	}).Info("Sending email notification")

	var provider EmailProvider

	// Select the email provider based on notification.Provider
	// switch notification.Provider {
	// case "sendgrid":
	// 	provider = &email_providers.SendGridProvider{}
	// 	// return email_providers.SendGrid(notification)
	// // case "another_email_provider":
	// // 	return email_providers.AnotherEmailProvider(notification)
	// default:
	// 	log.Println("Unsupported email provider:", notification.Provider)
	// 	return fmt.Errorf("unsupported email provider: %s", notification.Provider)
	// }

	return provider.Send(ctx, notification)
}

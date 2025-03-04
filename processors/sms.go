package processors

import (
	"context"
	"notify/config"
	"notify/models"
	"notify/repositories"
	"notify/services"
	"time"
)

// SMSProcessor handles processing of SMS notifications
func SMSProcessor(ctx context.Context, repo repositories.Repositories, notification models.Notification) error {
	var smsNotification models.SMSNotification

	smsNotification.CreatedAt = time.Now()
	smsNotification.UpdatedAt = time.Now()
	smsNotification.Recipient = notification.Recipient
	smsNotification.Status = "pending"
	smsNotification.NotificationID = notification.ID
	smsNotification.Message = notification.Payload["message"]
	smsNotification.Type = notification.Payload["type"]
	smsNotification.Provider = notification.Provider

	insertedID, err := repo.SmsRepo.SaveSMS(ctx, &smsNotification)
	if err != nil {
		config.Logger.WithError(err).Error("Error saving sms notification")
		return err
	}

	smsNotification.ID = insertedID.Hex()

	res, err := services.SendSMS(ctx, smsNotification)

	if err != nil {
		repo.SmsRepo.UpdateSMSResponse(ctx, smsNotification.ID, string(err.Error()), "error")
		return err
	}

	repo.SmsRepo.UpdateSMSResponse(ctx, smsNotification.ID, string(res), "success")

	return nil
}

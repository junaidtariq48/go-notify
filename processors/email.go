package processors

import (
	"context"
	"notify/config"
	"notify/models"
	"notify/repositories"
	"notify/services"
	"time"
)

func EmailProcessor(ctx context.Context, repo repositories.Repositories, notification models.Notification) error {
	var emailNotification models.EmailNotification

	emailNotification.CreatedAt = time.Now()
	emailNotification.UpdatedAt = time.Now()
	emailNotification.Recipient = notification.Recipient
	emailNotification.Status = "pending"
	emailNotification.NotificationID = notification.ID
	emailNotification.DynamicData = notification.Payload
	emailNotification.Type = notification.Type
	emailNotification.Provider = notification.Provider
	emailNotification.TemplateID = notification.TemplateId
	emailNotification.RecipientName = notification.RecipientName

	insertedID, err := repo.EmailRepo.SaveEmail(ctx, &emailNotification)
	if err != nil {
		config.Logger.WithError(err).Error("Error saving email notification")
		return err
	}

	emailNotification.ID = insertedID.Hex()

	res, err := services.SendEmail(ctx, emailNotification)

	if err != nil {
		repo.EmailRepo.UpdateEmailResposne(ctx, emailNotification.ID, string(err.Error()), "error")
		return err
	}

	repo.EmailRepo.UpdateEmailResposne(ctx, emailNotification.ID, string(res), "success")

	return nil
}

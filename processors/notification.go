package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"notify/config"
	"notify/models"
	"notify/repositories"
	"notify/services"

	"github.com/sirupsen/logrus"
)

func EmailProcessor(ctx context.Context, repo repositories.Repositories, notification models.Notification) error {

	var payload map[string]interface{}
	err := json.Unmarshal([]byte(notification.Payload), &payload)
	if err != nil {
		return err
	}

	dynamicData, ok := payload["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid 'data' field")
	}

	// Convert dynamicData to map[string]string
	dynamicTemplateData := make(map[string]string)
	for key, value := range dynamicData {
		if strValue, ok := value.(string); ok {
			dynamicTemplateData[key] = strValue
		} else {
			return fmt.Errorf("invalid value for key '%s'", key)
		}
	}

	var emailModel models.Email

	emailModel.To = payload["to"].(string)
	emailModel.From = config.AppConfig.FromEmail
	emailModel.Body = dynamicData
	emailModel.NotificaitonID = notification.ID
	emailModel.Status = "pending"

	// Save the notification to MongoDB
	insertedID, err := repo.EmailRepo.SaveEmail(ctx, &emailModel)

	if err != nil {

		config.Logger.WithFields(logrus.Fields{
			"type":         notification.Type,
			"notification": notification,
		}).Error("Error Processing notification")
	}

	fmt.Println("::EMAIL::", insertedID)
	return services.SendEmail(ctx, notification)

}

// SMSProcessor handles processing of SMS notifications
// func SMSProcessor(ctx context.Context, notification models.Notification) error {
// 	return services.SendSMS(ctx, notification)
// }

// You can add more processor functions here as needed
// For example:
// func PushNotificationProcessor(ctx context.Context, notification models.Notification) error {
//     return services.SendPushNotification(ctx, notification)
// }

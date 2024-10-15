package processors

import (
	"context"
	"notify/models"
	"notify/services"
)

func EmailProcessor(ctx context.Context, notification models.Notification) error {
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

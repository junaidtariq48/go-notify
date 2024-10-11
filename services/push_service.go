package services

import (
	"fmt"
	"log"
	"notify/config"
	"notify/models"
	"notify/services/push_providers"

	"github.com/sirupsen/logrus"
)

// PushProvider defines the interface for push notification providers
type PushProvider interface {
	Send(notification models.Notification) error
}

// // OneSignalProvider sends push notifications using OneSignal
// type OneSignalProvider struct{}

// func (o *OneSignalProvider) Send(notification models.Notification) error {
// 	fmt.Println("Sending Push Notification via OneSignal with payload:", notification.Payload)
// 	// TODO: Add OneSignal API logic here
// 	return nil
// }

// SendPushNotification sends a push notification using the appropriate provider
func SendPushNotification(notification models.Notification) error {
	// provider := &OneSignalProvider{} // We are using OneSignal for push notifications
	// return provider.Send(notification)
	var provider PushProvider

	config.Logger.WithFields(logrus.Fields{
		"notification_id": notification.ID,
		"provider":        notification.Provider,
	}).Info("Sending push notification")

	// Select the push provider based on notification.Provider
	switch notification.Provider {
	case "firebase":
		provider = &push_providers.FirebaeProvider{}
	case "onesignal":
		provider = &push_providers.OneSignalProvider{}
	default:
		log.Println("Unsupported push provider:", notification.Provider)
		return fmt.Errorf("unsupported push provider: %s", notification.Provider)
	}

	return provider.Send(notification)
}

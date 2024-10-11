package push_providers

import (
	"fmt"
	"notify/models"
)

// OneSignalProvider sends push notifications using OneSignal
type OneSignalProvider struct{}

func (o *OneSignalProvider) Send(notification models.Notification) error {
	fmt.Println("Sending Push Notification via OneSignal with payload:", notification.Payload)
	// TODO: Add OneSignal API logic here
	return nil
}

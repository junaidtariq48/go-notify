package email_providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"notify/config"
	"notify/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct{}

// SendGrid processes the email notification using the SendGrid provider
func (s *SendGridProvider) Send(ctx context.Context, notification models.EmailNotification) ([]byte, error) {
	sg := sendgrid.NewSendClient(config.AppConfig.SendGridApiKey)

	from := mail.NewEmail(config.AppConfig.FromName, config.AppConfig.FromEmail)

	subject := notification.Subject
	to := mail.NewEmail(notification.RecipientName, notification.Recipient)
	message := mail.NewV3MailInit(from, subject, to)
	message.SetTemplateID(notification.TemplateID)

	// Add substitution data to the dynamic template
	for key, value := range notification.DynamicData {
		message.Personalizations[0].SetDynamicTemplateData(key, value)
	}

	response, err := sg.Send(message)
	if err != nil || response.StatusCode != http.StatusAccepted {
		fmt.Println(err.Error())
		return nil, err
	}

	responseJSON, _ := json.Marshal(response)
	return responseJSON, nil
}

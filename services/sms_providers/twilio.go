package sms_providers

import (
	"encoding/json"
	"fmt"
	"log"
	"notify/config"
	"notify/models"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// TwilioProvider sends SMS using Twilio
type TwilioProvider struct{}

// Twilio processes the SMS notification using the Twilio provider
func (t *TwilioProvider) Send(notification models.SMSNotification) ([]byte, error) {
	sid := config.AppConfig.TwilioAccountSID
	token := config.AppConfig.TwilioAuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: sid,
		Password: token,
	})
	params := &verify.CreateVerificationParams{}
	params.SetCustomCode(notification.Message)
	params.SetChannel("sms")
	params.SetTo(notification.Recipient)

	resp, err := client.VerifyV2.CreateVerification(config.AppConfig.TwilioVerifySID,
		params)

	log.Printf("Twilio: Sending SMS to %s with message %s", notification.Recipient, notification.Message)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	} else {
		responseJSON, _ := json.Marshal(resp)
		return responseJSON, nil
	}
}

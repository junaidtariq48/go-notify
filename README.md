# go-notify

Notification Service

Email & SMS Notification:
// sample body

```
{
"type": "email", // mandatory value => email or sms
"provider": "sendgrid", // mandatory value => sendgrid - ForEmail : twilio - ForSMS
"recipient": "junaid.tariq48@gmail.com", // mandatory value => email or phone in case of sms
"subject": "this is test subject and optional", // optional only for email
"recipient_name": "Junaid Tariq", // mandatory value => name of the customer
"template_id" : "d-a9a09c52df3f460bb5691919c54a59bf", // mandatory for email only value => email template id
"payload": {
"customer_name": "Junaid Tariq",
"OTP": "232334",
"minutes": "10"
}
}
```

// email payload. below is the sample payload used to send email for otp. this is object of all the variables used in email and
// these are all variables not same all the time for every email separate.

```
{
"customer_name": "Junaid Tariq",
"OTP": "232334",
"minutes": "10"
}
```

// sms payload.

```
"payload": {
"type": "null", // mandatory for sms value => verify or null. verify to send the otp to user
"message": "hello this message is from test aqary." // in case of type = verify then send code in message. "23425"
}
```

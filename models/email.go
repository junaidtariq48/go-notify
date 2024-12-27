package models

type Email struct {
	ID             string            `json:"id" bson:"_id,omitempty"`
	NotificationID string            `json:"notification_id" bson:"notification_id"`
	Recipient      string            `json:"recipient" bson:"recipient"`       // Email address
	Subject        string            `json:"subject" bson:"subject"`           // Email subject
	TemplateID     string            `json:"template_id" bson:"template_id"`   // SendGrid dynamic template ID
	DynamicData    map[string]string `json:"dynamic_data" bson:"dynamic_data"` // Variables for the template
}

package queues

import (
	"context"
	"encoding/json"
	"log"

	"notify/models"

	"github.com/streadway/amqp"
)

func PublishNotification(ctx context.Context, channel *amqp.Channel, queueName string, notification models.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to queue %s: %v", queueName, err)
	}
	return err
}

func ConsumeNotification(ctx context.Context, channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	messages, err := channel.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register consumer on queue %s: %v", queueName, err)
		return nil, err
	}
	return messages, nil
}

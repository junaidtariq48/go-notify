package queues

import (
	"encoding/json"
	"log"
	"notify/models"

	"github.com/streadway/amqp"
)

func EnqueueRabbitMQ(ch *amqp.Channel, queueName string, notification models.Notification) error {
	message, _ := json.Marshal(notification)
	// Declare a queue
	_, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue: %s", err)
		return err
	}

	// Publish message to queue
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %s", err)
		return err
	}

	log.Printf("Message published to queue: %s", queueName)
	return nil
}

func DequeueRabbitMQ(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	// Consume messages from the queue
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to consume from queue: %s", err)
		return nil, err
	}

	return msgs, nil
}

package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"notify/config"
	"notify/models"

	"github.com/streadway/amqp"
)

// // Enqueue a notification to RabbitMQ
func EnqueueRabbitNotification(ctx context.Context, channel *amqp.Channel, queue string, notification models.Notification) error {
	// Attempt to marshal the notification into JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	_, errr := channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // args
	)

	if errr != nil {
		return err
	}

	// Try publishing the message, if an error occurs check if the channel is closed
	err = channel.Publish(
		"",    // Default exchange
		queue, // Routing key (queue name)
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        notificationJSON,
		},
	)
	if err != nil {
		// Log the error and attempt to reconnect
		log.Printf("Failed to publish message to RabbitMQ: %v", err)

		// If the error is due to a closed channel/connection, try reopening it
		if err == amqp.ErrClosed {
			log.Println("RabbitMQ channel is closed. Reconnecting...")
			// Reinitialize the connection and channel (if needed)
			newChannel, err := reinitializeRabbitMQChannel()
			if err != nil {
				return fmt.Errorf("failed to reinitialize RabbitMQ channel: %v", err)
			}
			channel = newChannel // Reassign the channel to the new one

			// Retry publishing the message with the new channel
			err = channel.Publish(
				"",    // Default exchange
				queue, // Routing key (queue name)
				false, // Mandatory
				false, // Immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        notificationJSON,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to publish after reconnecting: %v", err)
			}
		} else {
			return err
		}
	}

	return nil
}

// reinitializeRabbitMQChannel handles reconnecting to RabbitMQ and re-establishing a new channel.
func reinitializeRabbitMQChannel() (*amqp.Channel, error) {
	// Assuming you have the original RabbitMQ URI and connection setup
	conn, err := amqp.Dial(config.AppConfig.RabbitMQUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect to RabbitMQ: %v", err)
	}

	// Recreate the channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create a new RabbitMQ channel: %v", err)
	}

	return channel, nil
}

// // Dequeue a notification from RabbitMQ
func DequeueRabbitNotification(ctx context.Context, channel *amqp.Channel, queue string) (*models.Notification, error) {
	// if channel == nil {
	// inspect, err := channel.QueueInspect(queue)
	// if inspect.Messages <= 0 {
	// 	return nil, nil
	// }

	msgs, err := channel.Consume(
		queue, // Queue
		"",    // Consumer
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, err
	}

	// Listen for messages
	for msg := range msgs {
		var notification models.Notification
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			return nil, err
		}
		return &notification, nil
	}
	return nil, nil
}

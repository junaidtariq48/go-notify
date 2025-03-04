package amqp

import (
	"log"
	"notify/config"

	"github.com/streadway/amqp"
)

func InitRabbitMQ() *amqp.Connection {
	rabbitmqURI := config.AppConfig.RabbitMQUrl
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	log.Println("Hala wallaaa from....... ", rabbitmqURI)
	return conn
}

func DeclareQueue(channel *amqp.Channel, queueName string) {
	_, err := channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue %s: %v", queueName, err)
	}
}

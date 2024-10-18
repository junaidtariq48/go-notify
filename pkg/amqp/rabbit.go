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

	log.Println("Connected to RabbitMQ")
	return conn
}

package core

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitChannel(exchangeName string) (*amqp.Channel, error) {
	rabbitConnectionString := os.Getenv("RABBIT_CONNECTION_STRING")

	conn, err := amqp.Dial(rabbitConnectionString)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return nil, err
	}

	return ch, nil
}

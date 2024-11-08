package publisher

import (
	"context"
	"encoding/json"
	domain "pt/isep/insis/product-command/domain"
	"pt/isep/insis/product-command/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductPublisher interface {
	PublishMessage(ctx context.Context, message *domain.Event) *utils.Error
}

type productPublisherStruct struct {
	channel      *amqp.Channel
	exchangeName string
	routingKey   string
}

func NewProductPublisher(productChannel *amqp.Channel, exchangeName string, routingKey string) ProductPublisher {
	return &productPublisherStruct{productChannel, exchangeName, routingKey}
}

func (p *productPublisherStruct) PublishMessage(ctx context.Context, message *domain.Event) *utils.Error {

	bytes, err := json.Marshal(message)

	if err != nil {
		return &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	err = p.channel.PublishWithContext(ctx,
		p.exchangeName, // exchange
		p.routingKey,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})

	if err != nil {
		return &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	return nil
}

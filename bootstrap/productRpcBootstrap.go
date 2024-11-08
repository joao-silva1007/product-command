package bootstrap

import (
	"context"
	"encoding/json"
	"log"
	"pt/isep/insis/product-command/messages"
	"pt/isep/insis/product-command/service"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductRpcBootstrap interface {
	Bootstrap()
}

type ProductRpcBootstrapStruct struct {
	channel        *amqp.Channel
	productService service.ProductService
}

func NewProductRpcBootstrap(productChannel *amqp.Channel, productService service.ProductService) *ProductRpcBootstrapStruct {
	return &ProductRpcBootstrapStruct{productChannel, productService}
}

func (s *ProductRpcBootstrapStruct) Bootstrap() {
	q, err := s.channel.QueueDeclare(
		"product.rpc", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := s.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := uuid.NewString()

	go func() {
		for m := range msgs {
			if corrId == m.CorrelationId {
				ctx := context.TODO()
				event := new(messages.ProductRpcMessage)
				err := json.Unmarshal(m.Body, &event)
				if err != nil {
					return
				}

				for _, elem := range event.Products {
					s.productService.CreateProduct(ctx, elem)
				}

			}
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

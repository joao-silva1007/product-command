package bootstrap

import (
	"context"
	"encoding/json"
	"log"
	"pt/isep/insis/product-command/messages"
	"pt/isep/insis/product-command/service"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/mgo.v2/bson"
)

type ProductRpcBootstrap interface {
	Bootstrap()
}

type productRpcBootstrapStruct struct {
	channel        *amqp.Channel
	productService service.ProductService
}

func NewProductRpcBootstrap(productChannel *amqp.Channel, productService service.ProductService) ProductRpcBootstrap {
	return &productRpcBootstrapStruct{productChannel, productService}
}

func (s *productRpcBootstrapStruct) Bootstrap() {
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
	ctx := context.TODO()

	corrId := uuid.NewString()

	teste, _ := json.Marshal(bson.M{"teste": "teste"})

	err = s.channel.PublishWithContext(ctx,
		"",                     // exchange
		"product.rpc.requests", // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          teste,
		})
	failOnError(err, "Failed to publish a message")

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

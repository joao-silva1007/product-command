package listener

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"pt/isep/insis/product-command/domain"
	"pt/isep/insis/product-command/messages"
	"pt/isep/insis/product-command/service"
)

type ProductListener interface {
	StartListening()
}

type ProductListenerStruct struct {
	channel        *amqp.Channel
	exchangeName   string
	routingKey     string
	productService service.ProductService
}

func NewProductListener(productChannel *amqp.Channel, exchangeName string, routingKey string, productService service.ProductService) *ProductListenerStruct {
	return &ProductListenerStruct{productChannel, exchangeName, routingKey, productService}
}

func (s *ProductListenerStruct) StartListening() {
	q, err := s.channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = s.channel.QueueBind(
		q.Name,         // queue name
		"",             // routing key
		s.exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

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

	go func() {
		for m := range msgs {
			if m.RoutingKey != s.routingKey {
				ctx := context.TODO()
				event := new(domain.Event)
				err := json.Unmarshal(m.Body, &event)
				if err != nil {
					return
				}

				switch event.Type {
				case domain.CreateProduct:
					s.processCreateProductMessage(ctx, event.Body)
				case domain.UpdateProduct:
					s.processUpdateProductMessage(ctx, event.Body)
				case domain.DeleteProduct:
					s.processDeleteProductMessage(ctx, event.Body)
				}
			}
		}
	}()
}

func (s *ProductListenerStruct) processCreateProductMessage(ctx context.Context, bytes []byte) {
	createProductMessage := new(messages.CreateProductMessage)
	err := json.Unmarshal(bytes, &createProductMessage)
	if err != nil {
		return
	}

	_, err2 := s.productService.CreateProduct(ctx, createProductMessage.Product)
	log.Print(err2.BaseError)
}

func (s *ProductListenerStruct) processUpdateProductMessage(ctx context.Context, bytes []byte) {
	updateProductMessage := new(messages.UpdateProductMessage)
	err := json.Unmarshal(bytes, &updateProductMessage)
	if err != nil {
		return
	}

	_, err2 := s.productService.UpdateProductBySku(ctx, updateProductMessage.SKU, updateProductMessage.Product)
	log.Print(err2.BaseError)
}

func (s *ProductListenerStruct) processDeleteProductMessage(ctx context.Context, bytes []byte) {
	deleteProductMessage := new(messages.DeleteProductMessage)
	err := json.Unmarshal(bytes, &deleteProductMessage)
	if err != nil {
		return
	}

	err2 := s.productService.DeleteProductBySku(ctx, deleteProductMessage.SKU)
	log.Print(err2.BaseError)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

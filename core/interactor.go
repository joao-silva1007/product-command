package core

import (
	"pt/isep/insis/product-command/bootstrap"
	"pt/isep/insis/product-command/controller"
	"pt/isep/insis/product-command/listener"
	"pt/isep/insis/product-command/publisher"
	"pt/isep/insis/product-command/repository"
	"pt/isep/insis/product-command/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Interactor interface {
	NewAppHandler() AppHandler
	NewProductRepository() repository.ProductRepository
	NewProductService() service.ProductService
	NewProductController(productService service.ProductService) controller.ProductController
}

type InteractorStruct struct {
	databaseConn  *mongo.Database
	rabbitChannel *amqp.Channel
	routingKey    string
}

func NewInteractor(conn *mongo.Database, rabbitChannel *amqp.Channel, routingKey string) *InteractorStruct {
	return &InteractorStruct{conn, rabbitChannel, routingKey}
}

func (i *InteractorStruct) NewAppHandler() *AppHandlerStruct {
	appHandler := &AppHandlerStruct{}
	productService := i.NewProductService()
	appHandler.ProductController = i.NewProductController(productService)
	appHandler.ProductListener = i.NewProductListener(productService)
	productBootstrap := i.NewProductRpcBootstrapper(productService)
	productBootstrap.Bootstrap()
	appHandler.ProductListener.StartListening()

	return appHandler
}

func (i *InteractorStruct) NewProductRepository() repository.ProductRepository {
	return repository.NewProductRepository(i.databaseConn.Collection("products"))
}

func (i *InteractorStruct) NewEventRepository() repository.EventRepository {
	return repository.NewEventRepository(i.databaseConn.Collection("events"))
}

func (i *InteractorStruct) NewProductPublisher() publisher.ProductPublisher {
	return publisher.NewProductPublisher(i.rabbitChannel, "product-fanout", i.routingKey)
}

func (i *InteractorStruct) NewProductListener(productService service.ProductService) listener.ProductListener {
	return listener.NewProductListener(i.rabbitChannel, "product-fanout", i.routingKey, productService)
}

func (i *InteractorStruct) NewProductService() service.ProductService {
	return service.NewProductService(i.NewProductRepository(), i.NewEventRepository(), i.NewProductPublisher())
}

func (i *InteractorStruct) NewProductController(productService service.ProductService) controller.ProductController {
	return controller.NewProductController(productService)
}

func (i *InteractorStruct) NewProductRpcBootstrapper(productService service.ProductService) bootstrap.ProductRpcBootstrap {
	return bootstrap.NewProductRpcBootstrap(i.rabbitChannel, productService)
}

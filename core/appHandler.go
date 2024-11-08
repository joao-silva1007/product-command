package core

import (
	"pt/isep/insis/product-command/controller"
	"pt/isep/insis/product-command/listener"
)

type AppHandler interface {
	controller.ProductController
	listener.ProductListener
}

type appHandlerStruct struct {
	controller.ProductController
	listener.ProductListener
}

package routers

import (
	"pt/isep/insis/product-command/core"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(api fiber.Router, h core.AppHandler) {
	route := api.Group("/products")

	route.Post("/", h.CreateProduct)
	route.Patch("/:sku", h.UpdateProductBySku)
	route.Delete("/:sku", h.DeleteProductBySku)
}

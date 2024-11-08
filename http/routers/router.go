package routers

import (
	"pt/isep/insis/product-command/core"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, h core.AppHandler) {
	apiRoutes := app.Group("/")

	ProductRoutes(apiRoutes, h)
}

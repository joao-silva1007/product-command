package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type responseMessage struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func PageNotFound(c *fiber.Ctx, data interface{}) error {
	res := &responseMessage{
		Code:    fiber.StatusNotFound,
		Message: "Page Not Found",
		Data:    data,
	}
	return c.Status(fiber.StatusNotFound).JSON(res)
}

func ResponseWithBodyAndStatusCode(c *fiber.Ctx, data interface{}, statusCode int) error {
	return c.Status(statusCode).JSON(data)
}

func DatabaseResponseWithoutData(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func ErrorResponse(c *fiber.Ctx, data interface{}, statusCode int) error {
	res := &responseMessage{
		Code: fiber.StatusInternalServerError,
		Data: data,
	}
	return c.Status(statusCode).JSON(res)
}

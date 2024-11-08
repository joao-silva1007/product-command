package main

import (
	"github.com/google/uuid"
	"log"
	"os"
	"pt/isep/insis/product-command/core"
	"pt/isep/insis/product-command/helpers"
	"pt/isep/insis/product-command/http/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	databaseConn, databaseError := core.ConnectDatabase()
	if databaseError != nil {
		log.Fatal(databaseError)
		return
	}

	exchangeName := os.Getenv("PRODUCT_EXCHANGE_NAME")
	channel, rabbitChannelErr := core.SetupRabbitChannel(exchangeName)
	if rabbitChannelErr != nil {
		log.Fatal(rabbitChannelErr)
		return
	}

	interactor := core.NewInteractor(databaseConn, channel, uuid.NewString())
	handler := interactor.NewAppHandler()

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New())

	routers.NewRouter(app, handler)

	app.Use(func(c *fiber.Ctx) error {
		return helpers.PageNotFound(c, nil)
	})

	port := os.Getenv("PORT")

	initError := app.Listen(":" + port)
	if initError != nil {
		panic(initError)
	}
}

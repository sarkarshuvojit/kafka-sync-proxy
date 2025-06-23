package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarkarshuvojit/kafka-sync-proxy/internal/controller"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/v1")
	api.Post("/", controller.HandleRest)
}

func Start() {
	app := fiber.New(fiber.Config{
		AppName: "Kafka Sync Proxy",
	})
	setupRoutes(app)
	app.Listen(":8420")
}

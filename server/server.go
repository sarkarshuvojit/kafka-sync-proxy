package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarkarshuvojit/kafka-sync-proxy/controller"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/v1")
	api.Post("/", controller.HandleRest)
}

func Start() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(":8420")
}

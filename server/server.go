package server

import (
	"github.com/gofiber/fiber/v2"
	"shuvojit.in/asc/controller"
)




func setupRoutes(app *fiber.App) {
	api := app.Group("/v1")
	api.Post("/", controller.HandleRest)
}

func Start() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(":3000")
}

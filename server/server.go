package server

import (
	"github.com/gofiber/fiber/v2"
)

func testApi(c *fiber.Ctx) error {
    params := string(c.Request().URI().QueryString())
	c.Status(200).SendString(params)
    return nil
}
func Start() {
	app := fiber.New()
	app.Get("/", testApi)
	app.Listen(":3000")
}

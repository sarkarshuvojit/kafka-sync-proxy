package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"shuvojit.in/asc/messaging"
	"shuvojit.in/asc/messaging/kafka"
)

type Request struct {
    RequestTopic string `json:"requestTopic"`
    ResponseTopic string `json:"responseTopic"`
    Payload string `json:"payload"`
}

func requestResponseBlock(
    requestTopic string,
    responseTopic string,
    payload string,
) (string, error) {
    brokers := []string{ "localhost:29092" }
    var messagingProvider messaging.MessagingProvider
    messagingProvider = &kafka.Kafka{
        Brokers: brokers,
    }
    
    if err := messagingProvider.Send("GGG1155", requestTopic, []byte("{}")); err != nil {
        return "", errors.New("Could not send")
    }

    return "Kuch toh hua hai", nil
}

func handle(c *fiber.Ctx) error {
    request := Request{}
    if err := json.Unmarshal(c.Body(), &request); err != nil {
        c.Status(400).JSON(map[string]string{})
        return errors.New("Invalid request body")
    }
    log.Printf("Request: %v", request)
    requestResponseBlock(
        request.RequestTopic, 
        request.ResponseTopic, 
        request.Payload,
    )
    c.Status(200).JSON(map[string]string{
        "message": "fetched successfully",
    })
    return nil
}

func setupRoutes(app *fiber.App) {
    api := app.Group("/v1")
    api.Post("/", handle)
}

func Start() {
	app := fiber.New()
    setupRoutes(app)
	app.Listen(":3000")
}

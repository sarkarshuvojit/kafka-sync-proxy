package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"shuvojit.in/asc/messaging"
	"shuvojit.in/asc/messaging/kafka"
)

type Request struct {
    RequestTopic string `json:"requestTopic"`
    ResponseTopic string `json:"responseTopic"`
    Payload string `json:"payload"`
}

func createEventId() string {
    return uuid.NewString()
}

func requestResponseBlock(
    requestTopic string,
    responseTopic string,
    payload string,
) (string, error) {
    brokers := []string{ "localhost:29092" }
    var messagingProvider messaging.MessagingProvider
    messagingProvider = &kafka.Kafka{ Brokers: brokers }

    key := createEventId()
    
    if err := messagingProvider.Send(
        key,
        requestTopic, 
        []byte(payload),); 
        err != nil {
        return "", errors.New("Could not send")
    }

    return "Fetched response ", nil
}

func handle(c *fiber.Ctx) error {
    request := new(Request)
    if err := c.BodyParser(request); err != nil {
        c.Status(400).JSON(map[string]string{})
        return errors.New("Invalid request body")
    }
    log.Printf("Request: %v", request)
    res, err := requestResponseBlock(
        request.RequestTopic, 
        request.ResponseTopic, 
        request.Payload,
    )
    if err != nil {
        c.Status(400).JSON(map[string]string{})
        return errors.New("Did not get response")
    }
    c.Status(200).JSON(map[string]string{
        "message": "fetched successfully",
        "response": res,
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

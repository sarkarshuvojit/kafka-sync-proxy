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
) ([]byte, error) {

    brokers := []string{ "localhost:29092" }

    var messagingProvider messaging.MessagingProvider
    messagingProvider = &kafka.Kafka{ Brokers: brokers }

    msg, err := messagingProvider.SendAndReceive(
        requestTopic, 
        responseTopic,
        []byte(payload)); 

    if err != nil {
        log.Println("Error Receiving msg")
        return nil, err
    }

    return msg, nil
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
        return c.Status(400).JSON(map[string]string{
            "error": err.Error(),
        })
    }

    var response interface{}
    json.Unmarshal(res, &response)

    return c.Status(200).JSON(map[string]interface{}{
        "message": "fetched successfully",
        "response": response,
    })
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

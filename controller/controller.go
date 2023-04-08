package controller

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sarkarshuvojit/kafka-sync-proxy/messaging"
	"github.com/sarkarshuvojit/kafka-sync-proxy/messaging/kafka"
	"github.com/sarkarshuvojit/kafka-sync-proxy/service"
	"github.com/sarkarshuvojit/kafka-sync-proxy/types"
)

func HandleRest(c *fiber.Ctx) error {
	request := new(types.BlockingRequestDto)
	if err := c.BodyParser(request); err != nil {
		c.Status(400).JSON(map[string]string{})
		return errors.New("Invalid request body")
	}
	log.Printf("Request: %v", request)

	messagingProvider := &kafka.Kafka{Brokers: request.Brokers, Timeout: 5}
	blockingService := &service.BlockingService{Provider: messagingProvider}

	payloadAsBytes, _ := json.Marshal(request.Payload)

	res, err := blockingService.RequestResponseBlock(
		request.RequestTopic,
		request.ResponseTopic,
		string(payloadAsBytes),
	)
	if err != nil {
		status := 400
		if err == messaging.TimeoutErr {
			status = 408
		}
		return c.Status(status).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	var response interface{}
	json.Unmarshal(res, &response)

	return c.Status(200).JSON(map[string]interface{}{
		"response": response,
	})
}

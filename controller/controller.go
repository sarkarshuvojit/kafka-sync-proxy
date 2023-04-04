package controller

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"shuvojit.in/asc/service"
	"shuvojit.in/asc/types"
)

func HandleRest(c *fiber.Ctx) error {
	request := new(types.BlockingRequestDto)
	if err := c.BodyParser(request); err != nil {
		c.Status(400).JSON(map[string]string{})
		return errors.New("Invalid request body")
	}
	log.Printf("Request: %v", request)

	res, err := service.RequestResponseBlock(
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
		"message":  "fetched successfully",
		"response": response,
	})
}


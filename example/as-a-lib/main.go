package main

import (
	"fmt"

	"github.com/sarkarshuvojit/kafka-sync-proxy/pkg/messaging/kafka"
	"github.com/sarkarshuvojit/kafka-sync-proxy/pkg/service"
)

func main() {
	brokers := []string{"localhost:29092"}
	provider := kafka.Kafka{Brokers: brokers}
	service := service.BlockingService{Provider: provider}

	// The following line will
	// Push the Payload & Headers to the requestTopic
	// Wait for response in the responseTopic
	// Return with respone PayLoad & Headers
	//
	// By this we can achieve inter-service request reply pattern
	// while using a very simple interface
	response, err := service.RequestResponseBlock("create_user", "user_created", "{}", "{}")

	if err != nil {
		panic(err)
	}

	fmt.Println(response.Headers)
	fmt.Println(response.Payload)
}

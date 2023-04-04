package service

import (
	"log"

	"shuvojit.in/asc/messaging"
	"shuvojit.in/asc/messaging/kafka"
)
func RequestResponseBlock(
	requestTopic string,
	responseTopic string,
	payload string,
) ([]byte, error) {

	brokers := []string{"localhost:29092"}

	var messagingProvider messaging.MessagingProvider
	messagingProvider = &kafka.Kafka{Brokers: brokers}

	msg, err := messagingProvider.SendAndReceive(
		requestTopic,
		responseTopic,
		[]byte(payload))

	if err != nil {
		log.Println("Error Receiving msg")
		return nil, err
	}

	return msg, nil
}

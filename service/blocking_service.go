package service

import (
	"log"

	"github.com/sarkarshuvojit/kafka-sync-proxy/messaging"
)

type BlockingService struct {
	Provider messaging.MessagingProvider
}

func (b BlockingService) RequestResponseBlock(
	requestTopic string,
	responseTopic string,
	payload string,
) ([]byte, error) {

	messagingProvider := b.Provider

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

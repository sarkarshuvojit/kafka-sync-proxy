package service

import (
	"log"

	"github.com/sarkarshuvojit/kafka-sync-proxy/pkg/messaging"
)

type BlockingService struct {
	Provider messaging.MessagingProvider
}

// Creates new BlockingService with a messaging service provider
// Which will be used as the transport protocol
func NewBlockingService(provider messaging.MessagingProvider) *BlockingService {
	return &BlockingService{
		Provider: provider,
	}
}

// This function pushes the payload and headers to requestTopic,
// and waits for a message in responseTopic with the same key.
func (b BlockingService) RequestResponseBlock(
	requestTopic string,
	responseTopic string,
	payload string,
	headers string,
) (*messaging.SendAndReceiveResponse, error) {

	messagingProvider := b.Provider

	msg, err := messagingProvider.SendAndReceive(
		requestTopic,
		responseTopic,
		[]byte(payload),
		[]byte(headers))

	if err != nil {
		log.Println("Error Receiving msg")
		return nil, err
	}

	return msg, nil
}

package messaging

type MessagingProvider interface {
	SendAndReceive(
		requestTopic string,
		responseTopic string,
		payload []byte,
	) ([]byte, error)
}

package messaging

type SendAndReceiveResponse struct {
	Payload []byte
	Headers []byte
}

type MessagingProvider interface {
	SendAndReceive(
		requestTopic string,
		responseTopic string,
		payload []byte,
		headers []byte,
	) (*SendAndReceiveResponse, error)
}

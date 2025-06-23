package messaging

// SendAndReceiveResponse defines a generic structure for the response recieved from the response topic
type SendAndReceiveResponse struct {
	Payload []byte
	Headers []byte
}

// MessagingProvider is a common interface which provides a single blocking method called SendAndReceive
// This provider can be used in future to implement for different message providers
type MessagingProvider interface {
	SendAndReceive(
		requestTopic string,
		responseTopic string,
		payload []byte,
		headers []byte,
	) (*SendAndReceiveResponse, error)
}

package messaging

type MessagingProvider interface {
    Send(key string, requestTopic string, paylod []byte) error
    SendAndReceive(
        requestTopic string,
        responseTopic string,
        payload []byte,
    ) ([]byte, error)
}

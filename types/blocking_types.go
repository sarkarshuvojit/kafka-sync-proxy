package types

type BlockingRequestDto struct {
	RequestTopic  string      `json:"requestTopic"`
	ResponseTopic string      `json:"responseTopic"`
	Payload       interface{} `json:"payload"`
	Headers       interface{} `json:"headers"`
	Brokers       []string    `json:"brokers"`
}

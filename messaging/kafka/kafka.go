package kafka

import (
	"errors"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type Kafka struct {
    Brokers []string `json:"brokers"`
}

func (k Kafka) GetProducer() (sarama.SyncProducer, error) {
    log.Printf("Getting Producer for %v", k.Brokers)
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    conn, err := sarama.NewSyncProducer(k.Brokers, config)
    if err != nil {
        log.Println("Could not create producer")
        return nil, err
    }
    return conn, nil
}

func (k Kafka) Send(
    key string,
    requestTopic string,
    payload []byte,
) error {
    producer, err := k.GetProducer()
    if err != nil {
        log.Println("Error: %v", err)
        return errors.New(fmt.Sprintf("Could not produce to %v", k.Brokers))
    }
    partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
        Topic: requestTopic,
        Value: sarama.StringEncoder(string(payload)),
        Key: sarama.StringEncoder(key),
    })

    if err != nil {
        log.Println("Error: %v", err)
        return errors.New(fmt.Sprintf("Could not produce to %v", k.Brokers))
    }

    log.Printf("Produced to %s[%d, %d]", requestTopic, partition, offset)
    return nil
}


func (k Kafka) SendAndReceive(
    requestTopic string,
    payload []byte,
) ([]byte, error) {
    return []byte(""), nil
}
    

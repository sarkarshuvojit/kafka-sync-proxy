package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sarkarshuvojit/kafka-sync-proxy/pkg/messaging"
)

type Kafka struct {
	Brokers []string `json:"brokers,omitempty"`
	Timeout int      `json:"timeout,omitempty"`
}

// Creates new KafkaProvider with a set of brokers and timeout
// after which it will stop listening to the response channel
func NewKafkaProvider(brokers []string, timeout int) *Kafka {
	return &Kafka{
		Brokers: brokers,
		Timeout: timeout,
	}
}

func (k Kafka) createEventId() string {
	return uuid.NewString()
}

func (k Kafka) getProducer() (sarama.SyncProducer, error) {
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

func (k Kafka) getConsumer() (sarama.Consumer, error) {
	log.Printf("Getting Producer for %v", k.Brokers)
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	conn, err := sarama.NewConsumer(k.Brokers, config)
	if err != nil {
		log.Println("Could not create producer")
		return nil, err
	}
	return conn, nil
}

func (k Kafka) send(
	key string,
	requestTopic string,
	payload []byte,
	headers []byte,
) error {
	producer, err := k.getProducer()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return errors.New(fmt.Sprintf("Could not produce to %v", k.Brokers))
	}
	defer producer.Close()

	var recordHeaders []sarama.RecordHeader
	var headersMap map[string]string

	if err := json.Unmarshal(headers, &headersMap); err != nil {
		return err
	}

	for k, v := range headersMap {
		recordHeaders = append(recordHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:   requestTopic,
		Value:   sarama.StringEncoder(string(payload)),
		Key:     sarama.StringEncoder(key),
		Headers: recordHeaders,
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return errors.New(fmt.Sprintf("Could not produce to %v", k.Brokers))
	}

	log.Printf("Produced to %s[%d, %d]", requestTopic, partition, offset)
	return nil
}

func (k Kafka) receive(
	ctx context.Context,
	topic string,
	coorelationId string,
	timeoutInSeconds int,
) (*messaging.SendAndReceiveResponse, error) {
	worker, err := k.getConsumer()
	if err != nil {
		log.Println("Could not create consumer")
		return nil, err
	}
	defer worker.Close()

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Println("Cound not start consumer")
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		time.Second*time.Duration(timeoutInSeconds),
	)
	defer cancel()

	// Get signal for finish
	responseFoundCh := make(chan *messaging.SendAndReceiveResponse)
	responseFoundErrCh := make(chan error)

	go func(ctx context.Context) {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
				responseFoundErrCh <- err
				return

			case msg := <-consumer.Messages():
				if string(msg.Key) == coorelationId {
					// Todo: Dynamically set logic for finding reuqest pairs
					fmt.Println("Found coorelationId")
					headersMap := map[string]string{}
					for _, recordHeader := range msg.Headers {
						log.Printf("recordHeader: %s", recordHeader.Key)
						log.Printf("recordHeader: %s", recordHeader.Value)
						headersMap[string(recordHeader.Key)] = string(recordHeader.Value)
					}
					headersAsBytes, _ := json.Marshal(headersMap)
					responseFoundCh <- &messaging.SendAndReceiveResponse{
						Payload: msg.Value,
						Headers: headersAsBytes,
					}
					return
				}

			case <-ctx.Done():
				fmt.Println("Context cancelled")
				consumer.Close()
				responseFoundErrCh <- errors.New("Context cancelled")
				return
			}
		}
	}(ctxTimeout)

	select {
	case <-ctxTimeout.Done():
		fmt.Printf("Context cancelled: %v\n", ctxTimeout.Err())
		return nil, messaging.TimeoutErr

	case err := <-responseFoundErrCh:
		fmt.Println("Unkown error")
		return nil, err

	case response := <-responseFoundCh:
		return response, nil
	}

}

// SendAndReceive accepts request and reponse topic
// Pushes a message in request topic and expects the response in the response topic with the same key
func (k Kafka) SendAndReceive(
	ctx context.Context,
	requestTopic string,
	responseTopic string,
	payload []byte,
	headers []byte,
) (*messaging.SendAndReceiveResponse, error) {

	coorelationId := k.createEventId()

	if err := k.send(coorelationId, requestTopic, payload, headers); err != nil {
		log.Println("Error sending message")
		return nil, err
	}

	res, err := k.receive(ctx, responseTopic, coorelationId, k.Timeout)
	if err != nil {
		return nil, err
	}

	return res, nil
}

package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type Kafka struct {
    Brokers []string `json:"brokers"`
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
) error {
    producer, err := k.getProducer()
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

func (k Kafka) receive(
    topic string,
    coorelationId string,
    timeoutInSeconds int,
) ([]byte, error) {
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

    sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

    ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutInSeconds))
	defer cancel()

	// Get signal for finish
	responseFoundCh := make(chan []byte)
	responseFoundErrCh := make(chan error)

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
                responseFoundErrCh <- err

			case msg := <-consumer.Messages():
                if string(msg.Key) == coorelationId {
                    fmt.Println("Found coorelationId")
                    responseFoundCh <- msg.Value
                }

			case <-sigchan:
				fmt.Println("Interrupt is detected")
				responseFoundErrCh <- errors.New("Interruped")
			}
		}
	}()

    select {
        case <- ctxTimeout.Done():
            fmt.Printf("Context cancelled: %v\n", ctxTimeout.Err())
            return nil, errors.New("Timeout error")

        case err := <- responseFoundErrCh:
            fmt.Println("Unkown error")
            return nil, err

        case response := <- responseFoundCh:
            return response, nil
    }

}

func (k Kafka) SendAndReceive(
    requestTopic string,
    responseTopic string,
    payload []byte,
) ([]byte, error) {

    coorelationId := k.createEventId()

    if err := k.send(coorelationId, requestTopic, []byte("{}")); err != nil {
        log.Println("Error sending message")
        return nil, err
    }

    res, err := k.receive(responseTopic, coorelationId, 5)
    if err != nil {
        return nil, err
    }

    return res, nil
}


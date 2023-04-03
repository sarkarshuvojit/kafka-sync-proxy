package kafka

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
)

type Kafka struct {
    Brokers []string `json:"brokers"`
}

func createEventId() string {
    return uuid.NewString()
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

func (k Kafka) GetConsumer() (sarama.Consumer, error) {
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

func (k Kafka) SendAndReceive(
    requestTopic string,
    responseTopic string,
    payload []byte,
) ([]byte, error) {
    log.Println("Send&Receive")

    coorelationId := createEventId()

    if err := k.Send(coorelationId, responseTopic, []byte("{}")); err != nil {
        log.Fatalln("Error sending message")
    }

    worker, err := k.GetConsumer()
    if err != nil {
        log.Fatalln("Could not create consumer")
    }

    consumer, err := worker.ConsumePartition(responseTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
    sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) | Key(%s)\n", msgCount, string(msg.Topic), string(msg.Value), string(msg.Key))
                if msgCount > 5 {
                    doneCh <- struct{}{}
                }
                if string(msg.Key) == coorelationId {
                    fmt.Println("Found coorelationId")
                    doneCh <- struct{}{}
                }
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}

    return []byte(""), nil
}

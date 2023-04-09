package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// Set up Kafka configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Define the topic to consume from and the topic to produce to
	consumeTopic := "testing.in"
	produceTopic := "testing.out"

	// Set up Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:29092"}, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			cancel()
		}
	}()

	// Set up Kafka producer
	producer, err := sarama.NewAsyncProducer([]string{"localhost:29092"}, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			cancel()
		}
	}()

	// Handle OS interrupts to gracefully shutdown consumer and producer
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		if err := consumer.Close(); err != nil {
			cancel()
		}
		if err := producer.Close(); err != nil {
			cancel()
		}

		cancel()
	}()

	// Set up Kafka consumer to listen to the consumeTopic
	consumerPartitionList, err := consumer.Partitions(consumeTopic)
	if err != nil {
		panic(err)
	}

	for partition := range consumerPartitionList {
		pc, err := consumer.ConsumePartition(consumeTopic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := pc.Close(); err != nil {
				panic(err)
			}
		}()
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Received message with key=%s, value=%s, topic=%s, partition=%d, offset=%d\n", string(msg.Key), string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

				var msgHeaders []sarama.RecordHeader
				for _, h := range msg.Headers {
					msgHeaders = append(msgHeaders, sarama.RecordHeader{
						Key:   h.Key,
						Value: h.Value,
					})
				}

				msgHeaders = append(msgHeaders, sarama.RecordHeader{
					Key:   []byte("newHeaderKey"),
					Value: []byte("header val"),
				})

				// Produce a new message to the produceTopic with the same key and payload
				producer.Input() <- &sarama.ProducerMessage{
					Topic:   produceTopic,
					Key:     sarama.StringEncoder(msg.Key),
					Value:   sarama.StringEncoder(msg.Value),
					Headers: msgHeaders,
				}
			}
		}(pc)
		fmt.Println("Started consumer")
	}

	for {
		select {
		case <-signals:
		case <-ctx.Done():
		}
	}
}

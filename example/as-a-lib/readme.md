# Using as a library

## Use case

Let's say in while you're developing some api there comes a need to produce a message to a topic, but for this specific case, your code also needs some confirmation which comes to a different topic altogether.

We can argue that the response is better sent later using sse/ws, but let's assume there's something that's blocking us from doing that. 

In that case, it will be very easy to use this as a client to consume respone from a request-response pattern.

## Tutorial

### Initialising the blocking service

```go
// List of brokers is required to initialise our kafka provider
brokers := []string{"localhost:29092"}

// Later we will have multiple provider types like redis/rabbitmq
provider := kafka.Kafka{Brokers: brokers}

// finally create the service
service := service.BlockingService{Provider: provider}
```

### Initiate request-reply operation

```go
response, err := service.RequestResponseBlock("requestTopic", "responseTopic", "{}", "{}")
```

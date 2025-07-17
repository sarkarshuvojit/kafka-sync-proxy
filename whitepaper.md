# Kafka Synchronous Proxy: Technical Whitepaper

## Abstract

The Kafka Synchronous Proxy is a web service that bridges the gap between traditional synchronous REST APIs and asynchronous Kafka message streaming. It enables request-response patterns over Kafka by implementing a correlation-based message routing system with timeout handling. This proxy allows systems to interact with Kafka-based microservices using familiar HTTP semantics while maintaining the benefits of message-driven architecture.

## 1. Introduction

Modern distributed systems often employ event-driven architectures using Apache Kafka for scalability and decoupling. However, many client applications and legacy systems still require synchronous request-response patterns. The Kafka Synchronous Proxy addresses this architectural mismatch by providing a REST API that translates synchronous HTTP requests into asynchronous Kafka message exchanges.

## 2. System Architecture

### 2.1 Core Components

The system consists of four main architectural layers:

1. **HTTP Server Layer** (`internal/server/server.go`): Built on Fiber v2.43.0, handles HTTP requests and provides graceful shutdown capabilities
2. **Controller Layer** (`internal/controller/controller.go`): Manages request parsing, service orchestration, and response formatting
3. **Service Layer** (`pkg/service/blocking_service.go`): Implements the core business logic for synchronous message exchange
4. **Messaging Layer** (`pkg/messaging/kafka/kafka.go`): Provides Kafka-specific implementation of the messaging interface

### 2.2 Request Flow Architecture

```
HTTP Client → REST Endpoint → Controller → Blocking Service → Kafka Producer
                                                                    ↓
HTTP Response ← JSON Formatter ← Service Layer ← Kafka Consumer ← Response Topic
```

## 3. Technical Implementation

### 3.1 Correlation-Based Message Routing

The system implements a correlation ID mechanism to match requests with responses:

- **Correlation ID Generation**: Uses UUID v4 for unique message identification
- **Message Key Mapping**: Kafka message keys serve as correlation identifiers
- **Response Matching**: Consumer filters messages by correlation ID to find corresponding responses

### 3.2 Request Processing Flow

1. **Request Ingestion**: HTTP POST requests are received at `/v1/` endpoint
2. **Payload Parsing**: Request body is parsed into `BlockingRequestDto` structure containing:
   - `requestTopic`: Target Kafka topic for the request
   - `responseTopic`: Topic to monitor for responses
   - `payload`: Message payload (JSON)
   - `headers`: Custom headers to forward
   - `brokers`: Kafka broker addresses

3. **Message Production**: 
   - Generates unique correlation ID
   - Serializes payload and headers to JSON
   - Produces message to request topic with correlation ID as key

4. **Response Consumption**:
   - Starts consumer on response topic
   - Filters messages by correlation ID
   - Implements configurable timeout (default: 5 seconds)

### 3.3 Timeout Management

The system implements a robust timeout mechanism:

- **Configuration**: Timeout is configurable via `KSP_TIMEOUT` environment variable
- **Default Value**: 5 seconds if not specified
- **Context Cancellation**: Uses Go context for graceful timeout handling
- **HTTP Status Mapping**: Returns HTTP 408 (Request Timeout) on timeout

### 3.4 Error Handling

Comprehensive error handling covers:

- **Kafka Connection Errors**: Producer/consumer initialization failures
- **Message Parsing Errors**: Invalid JSON in request/response
- **Timeout Errors**: No response received within timeout period
- **Context Cancellation**: Graceful shutdown scenarios

## 4. Message Format Specification

### 4.1 Request Message Format

**HTTP Request:**
```json
{
  "requestTopic": "com.example.service.request",
  "responseTopic": "com.example.service.response",
  "payload": {
    "data": "example payload"
  },
  "headers": {
    "authorization": "Bearer token"
  },
  "brokers": ["localhost:9092"]
}
```

**Kafka Message:**
- **Key**: UUID v4 correlation ID
- **Value**: JSON serialized payload
- **Headers**: Key-value pairs from request headers

### 4.2 Response Message Format

**Kafka Response:**
- **Key**: Same correlation ID from request
- **Value**: Service response payload
- **Headers**: Response headers from processing service

**HTTP Response:**
```json
{
  "response": {
    "result": "processed data"
  },
  "headers": {
    "authorization": "Bearer token",
    "custom-header": "value"
  }
}
```

## 5. Performance Characteristics

### 5.1 Scalability Features

- **Stateless Design**: Each request is independent, enabling horizontal scaling
- **Connection Pooling**: Kafka producers and consumers are created per request
- **Graceful Shutdown**: 30-second graceful shutdown timeout for request completion
- **Memory Efficiency**: Minimal memory footprint with automatic resource cleanup

### 5.2 Throughput Considerations

- **Synchronous Blocking**: Each request blocks until response or timeout
- **Concurrency**: Multiple requests handled concurrently by Fiber framework
- **Kafka Partition Strategy**: Currently uses partition 0 for consumption (single partition)

### 5.3 Latency Factors

- **Network Latency**: Kafka broker communication overhead
- **Processing Time**: Downstream service processing duration
- **Timeout Boundary**: Maximum 5-second response time (configurable)

## 6. Deployment and Configuration

### 6.1 Deployment Options

The service supports multiple deployment modes:

1. **Standalone Binary**: Direct execution with Go runtime
2. **Docker Container**: Containerized deployment with port 8420 exposure
3. **Docker Compose**: Multi-service orchestration support

### 6.2 Configuration Parameters

- **Port**: Fixed at 8420 (HTTP server)
- **Timeout**: Configurable via `KSP_TIMEOUT` environment variable
- **Kafka Brokers**: Specified per request (dynamic configuration)
- **Log Level**: Standard Go log output

## 7. Security Considerations

### 7.1 Authentication and Authorization

- **Header Passthrough**: Custom authentication headers forwarded to Kafka
- **Broker Security**: Supports standard Kafka security protocols
- **No Built-in Authentication**: Relies on external authentication mechanisms

### 7.2 Data Protection

- **Payload Encryption**: Supports encrypted payloads (transparent handling)
- **Header Sanitization**: Headers processed without modification
- **Memory Safety**: Automatic cleanup of sensitive data

## 8. Use Cases and Applications

### 8.1 Legacy System Integration

- **API Gateway Pattern**: Expose Kafka-based services as REST APIs
- **Microservice Migration**: Gradual migration from synchronous to asynchronous patterns
- **Protocol Translation**: Bridge HTTP and Kafka ecosystems

### 8.2 Event-Driven Architecture

- **Command-Query Separation**: Separate command and query operations
- **Service Orchestration**: Coordinate multiple service interactions
- **Request-Response Patterns**: Implement synchronous behavior in event-driven systems

## 9. Limitations and Considerations

### 9.1 Current Limitations

- **Single Partition Consumption**: Only consumes from partition 0
- **No Message Persistence**: No retry mechanism for failed messages
- **Fixed Timeout**: Global timeout applies to all requests
- **No Load Balancing**: No built-in consumer group management

### 9.2 Operational Considerations

- **Resource Usage**: Each request creates new Kafka connections
- **Error Recovery**: Manual intervention required for broker failures
- **Monitoring**: Limited built-in metrics and health checks

## 10. Future Enhancements

### 10.1 Proposed Improvements

- **Multi-Partition Support**: Distribute load across multiple partitions
- **Connection Pooling**: Reuse Kafka connections for better performance
- **Health Checks**: Implement health and readiness endpoints
- **Metrics Integration**: Add Prometheus/OpenTelemetry support

### 10.2 Architectural Evolution

- **Message Persistence**: Add message replay capabilities
- **Circuit Breaker**: Implement fault tolerance patterns
- **Dynamic Configuration**: Runtime configuration updates
- **Consumer Groups**: Support for scalable consumer patterns

## 11. Conclusion

The Kafka Synchronous Proxy provides a practical solution for integrating synchronous client applications with asynchronous Kafka-based services. Its correlation-based message routing, timeout handling, and stateless design make it suitable for production environments where request-response patterns are required over Kafka infrastructure.

The system's modular architecture and clear separation of concerns enable easy extension and customization for specific use cases. While current limitations exist, the foundation provides a solid base for future enhancements and enterprise-grade features.

## References

- **Apache Kafka**: https://kafka.apache.org/
- **Fiber Web Framework**: https://gofiber.io/
- **Shopify Sarama**: https://github.com/Shopify/sarama
- **Go Context Package**: https://pkg.go.dev/context
# Usage as a proxy

## Async microservice

If you have a microservice which does not support rest, but rather just listens to topics and responds to other topics like the example file `main.go`. 


In the main.go provided as an example, we are listening on `testing.in` while sending the response back to `testing.out`.

```bash
$ go run main.go
```

Then we can run kafka-sync-proxy as a service using the binary or docker. 

```
$ go install github.com/sarkarshuvojit/kafka-sync-proxy
$ kafka-sync-proxy
```

Then you can fire requests using a curl similar to the following:

```bash
curl --location --request POST 'http://localhost:8420/v1/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "requestTopic": "testing.in",
    "responseTopic": "testing.out",
    "payload": {
        "message": "testing the proxy",
        "nested": {
            "child": {
                "subChild": 42
            }
        }
    },
    "headers": {
        "customAuthHeader": "ggg"
    },
    "brokers": [
        "localhost:29092"
    ]
}'
```



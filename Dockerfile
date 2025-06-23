FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kafka-sync-proxy cmd/serve/main.go

EXPOSE 8420

CMD ["/app/kafka-sync-proxy"]

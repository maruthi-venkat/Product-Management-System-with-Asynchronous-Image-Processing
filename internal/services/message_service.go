package services

import (
    "github.com/streadway/amqp"
)

// MessageService handles message publishing to a queue (e.g., RabbitMQ)
type MessageService struct {
    Connection *amqp.Connection
}

// NewMessageService creates a new MessageService instance
func NewMessageService() (*MessageService, error) {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, err
    }

    return &MessageService{Connection: conn}, nil
}

// PublishMessage publishes a message to a RabbitMQ queue
func (ms *MessageService) PublishMessage(queueName, message string) error {
    ch, err := ms.Connection.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    err = ch.ExchangeDeclare(
        "direct_logs",   // name of the exchange
        "direct",        // type of exchange
        true,            // durable
        false,           // auto-deleted
        false,           // internal
        false,           // no-wait
        nil,             // arguments
    )
    if err != nil {
        return err
    }

    err = ch.Publish(
        "direct_logs",   // exchange
        queueName,       // routing key
        false,           // mandatory
        false,           // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    return err
}

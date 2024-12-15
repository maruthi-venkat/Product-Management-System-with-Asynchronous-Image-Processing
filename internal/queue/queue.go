package queue

import (
    "log"
    "github.com/streadway/amqp"
)

// PublishMessage publishes a message to the given queue.
func PublishMessage(queueName, message string) error {
    // Establish connection to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
        return err
    }
    defer conn.Close()

    // Create a channel
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return err
    }
    defer ch.Close()

    // Declare a queue
    q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
        return err
    }

    // Log queue declaration
    log.Printf("Queue declared: %s", q.Name)

    // Publish the message
    err = ch.Publish("", q.Name, false, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(message),
    })
    if err != nil {
        log.Fatalf("Failed to publish a message: %v", err)
        return err
    }

    log.Printf("Message successfully published to queue: %s", q.Name)

    return nil
}

package queue

import (
    "log"
    "github.com/streadway/amqp"
    "product-management/internal/services"  // Assuming image_processor is in the 'services' package
)

func StartWorker(queueName string) {
    // Establish connection to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Create a channel
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // Declare the queue
    _, err = ch.QueueDeclare(
        queueName,  // Queue name
        true,       // Durable: Queue survives server restart
        false,      // AutoDelete: Queue will not be deleted when no consumers
        false,      // Exclusive: Only this connection can access the queue
        false,      // NoWait: Don't wait for server confirmation
        nil,        // Arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare the queue: %v", err)
    }

    // Start consuming messages
    msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to start consuming from the queue: %v", err)
    }

    // Process each message
    for msg := range msgs {
        log.Printf("Received image URL: %s", string(msg.Body))

        // Process image and handle errors
        err := services.ProcessImage(string(msg.Body))
        if err != nil {
            log.Printf("Failed to process image: %v", err)
            ch.Nack(msg.DeliveryTag, false, true) // Retry the message if processing failed
        } else {
            // Acknowledge that the message was processed successfully
            ch.Ack(msg.DeliveryTag, false)
        }
    }
}

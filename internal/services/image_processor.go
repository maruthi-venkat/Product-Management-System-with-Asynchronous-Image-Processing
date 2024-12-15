package services

import (
    "fmt"
    "log"
    "os"
    "net/http"
    "image"
    "github.com/streadway/amqp"
    "github.com/nfnt/resize"
    "image/jpeg"
    "bytes"
    "io/ioutil"
    "path/filepath"
    "time"
    "database/sql"
   _ "github.com/lib/pq" 
)

// ProcessImageWorker listens to the queue and processes images asynchronously
func ProcessImageWorker(queueName string) {
    // Connect to RabbitMQ
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

    // Declare the queue to ensure it exists
    _, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }
    log.Printf("Queue declared: %s", queueName)

    // Consume messages from the queue
    msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to start consuming from the queue: %v", err)
    }
    log.Printf("Waiting for messages in queue: %s", queueName)

    // Process each message
    for msg := range msgs {
        log.Printf("Received message: %s", msg.Body)

        imageURL := string(msg.Body)
        err := ProcessImage(imageURL)
        if err != nil {
            log.Printf("Failed to process image: %v", err)
            // Negative Acknowledge the message if processing fails
            ch.Nack(msg.DeliveryTag, false, true)
            continue
        }

        // Acknowledge the message after successful processing
        err = ch.Ack(msg.DeliveryTag, false)
        if err != nil {
            log.Printf("Failed to acknowledge message: %v", err)
        } else {
            log.Printf("Message acknowledged: %s", msg.Body)
        }
    }
}


// ProcessImage handles the image processing (download, resize, compress, and save)
func ProcessImage(imageURL string) error {
    // Download the image
    resp, err := http.Get(imageURL)
    if err != nil {
        return fmt.Errorf("failed to download image: %v", err)
    }
    defer resp.Body.Close()

    // Read the image data
    imgData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read image data: %v", err)
    }

    // Decode the image
    img, _, err := image.Decode(bytes.NewReader(imgData))
    if err != nil {
        return fmt.Errorf("failed to decode image: %v", err)
    }

    // Resize the image (compression)
    resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

    // Create a directory to store processed images (if it doesn't exist)
    outputDir := "processed_images"
    err = os.MkdirAll(outputDir, 0755) // Create the directory if not exists
    if err != nil {
        return fmt.Errorf("failed to create directory: %v", err)
    }

    // Generate a unique filename using timestamp and image URL hash (or any unique approach)
    fileName := fmt.Sprintf("%s_%d.jpg", filepath.Base(imageURL), time.Now().Unix())

    // Create the output file in the directory
    outFile, err := os.Create(filepath.Join(outputDir, fileName))
    if err != nil {
        return fmt.Errorf("failed to create processed image: %v", err)
    }
    defer outFile.Close()

    // Compress and save as JPEG
    err = jpeg.Encode(outFile, resizedImg, nil)
    if err != nil {
        return fmt.Errorf("failed to encode image as JPEG: %v", err)
    }

    log.Printf("Image successfully processed and saved: %s", fileName)
    
    // After processing, update the database with the new image URL or file path
    productID := 2 // You need to pass the actual product ID
    err = UpdateProductImageInDatabase(productID, filepath.Join(outputDir, fileName))
    if err != nil {
        return fmt.Errorf("failed to update product image in database: %v", err)
    }
    
    return nil
}
// UpdateProductImageInDatabase updates the product's image URL in the database
func UpdateProductImageInDatabase(productID int, compressedImagePath string) error {
    // Replace with your actual database connection details
    db, err := sql.Open("postgres", "user=maruthi password=12345678 dbname=product_management sslmode=disable")
    if err != nil {
        return fmt.Errorf("failed to connect to the database: %v", err)
    }
    defer db.Close()
    // Normalize the file path for cross-platform compatibility
    compressedImagePath := filepath.Join(outputDir, fileName)
    compressedImagePath = filepath.ToSlash(compressedImagePath) // Convert to use forward slashes

    log.Printf("Updating product %d with image path: %s", productID, compressedImagePath)
    // Query to update the product record with the compressed image path
    query := "UPDATE products SET compressed_product_images = $1 WHERE id = $2"
    _, err = db.Exec(query, compressedImagePath, productID)
    if err != nil {
        return fmt.Errorf("failed to update the product record: %v", err)
    }

    log.Printf("Product image updated for product ID %d with compressed image: %s", productID, compressedImagePath)
    return nil
}
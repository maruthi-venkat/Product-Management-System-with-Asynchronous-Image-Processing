package main

import (
	"fmt"
	"log"
    "net/http"
	"os"
	"os/signal"
	"product-management/internal/api/handlers"
	"product-management/internal/db"
	"product-management/internal/services"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"product-management/internal/queue" // Import the queue package where StartWorker is defined
)

var ctx = context.Background()

func main() {
	// Database connection string
	connString := "postgres://maruthi:12345678@localhost:5432/product_management?sslmode=disable"

	// Initialize the database
	if err := db.InitDB(connString); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Initialize Redis client and cache service
	services.InitializeCache()
	if err := services.RedisClient.Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")

	// Initialize MessageService
	messageService, err := services.NewMessageService()
	if err != nil {
		log.Fatalf("Failed to initialize MessageService: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Routes for handling products
	log.Println("Server is starting...")
	r.POST("/products", func(c *gin.Context) {
		log.Println("Received a POST request to /products")
		handlers.CreateProduct(c, *messageService) // Pass MessageService here
	})

	r.GET("/products/:id", func(c *gin.Context) {
		handlers.GetProductByID(c, services.RedisClient) // Pass CacheService here
	})

	r.GET("/products", func(c *gin.Context) {
		handlers.GetProductsByUser(c, services.RedisClient) // Pass CacheService here
	})

	// Start the image processing worker in a separate goroutine
	go func() {
		log.Println("Starting image processing worker...")
		queue.StartWorker("image_processing_queue") // Provide the queue name
	}()

	// Create an HTTP server instance
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    // Graceful shutdown handling
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe error: %v", err)
        }
    }()

    // Setup signal handling to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit // Wait for an interrupt signal

    log.Println("Shutting down the server...")
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server Shutdown Failed: %v", err)
    }
    log.Println("Server exited gracefully")
}

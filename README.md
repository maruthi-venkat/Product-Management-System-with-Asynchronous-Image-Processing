# Product Image Processing System

A scalable backend system that processes and updates product images asynchronously. The system is designed to handle image processing tasks in parallel, update product records in a PostgreSQL database, and ensure high performance using RabbitMQ for message queuing.

## Features
- Asynchronous image processing with Go.
- Image compression and saving in the required format.
- Integration with RabbitMQ for task queuing.
- PostgreSQL database integration for storing product image URLs.

## Getting Started

### Prerequisites
- Go programming language (1.x)
- PostgreSQL database
- RabbitMQ server
- redis
- docker

### Installing

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/product-image-processing.git
   cd product-image-processing
   ```

2. Set up your PostgreSQL database:
   - Create a database and update the database connection details in the `.env` file.

3. Install the necessary dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

### Usage
- To trigger image processing, send a POST request to the `/process-image` endpoint with a product image URL.
  - Example:
    ```bash
    curl -X POST http://localhost:8080/process-image -d '{"image_url": "https://example.com/image.jpg"}'
    ```

### Testing
- The system includes unit tests for the image processing logic. You can run the tests using:
  ```bash
  go test ./...
  ```
  or else use postman

## Technologies Used
- **Go**: The programming language used for building the backend system.
- **PostgreSQL**: The database used for storing product image records.
- **RabbitMQ**: The message queuing system used for asynchronous image processing.
- **Image Processing Libraries**: Libraries used for compressing and saving images in the required format.

## Challenges
- Handling image URLs and ensuring proper error handling.
- Managing asynchronous tasks efficiently with RabbitMQ.

## Future Improvements
- Scaling the system to handle more complex image processing tasks.
- Implementing robust logging and monitoring.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

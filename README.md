# Ticket Processor

Ticket Processor is a Go-based application for processing receipts and managing points. This README provides an overview of the project structure and how to access the Swagger documentation for the API.

## Project Structure

- `cmd/`: Contains the main application entry point.
- `internal/`: Contains the core application code.
    - `api/`: Contains the API-related code.
        - `handlers/`: Contains the HTTP handlers for the API endpoints.
        - `middleware/`: Contains the middleware for the API.
    - `config/`: Contains the configuration loading and management code.
    - `services/`: Contains the business logic and service layer.
    - `storage/`: Contains the storage layer for data persistence.
    - `receipt/`: Implements methods for calculating points based on receipt data.
    - `validation/`: Contains the validation logic for the application.
- `pkg/`: Contains shared packages used across the application.
    - `logger/`: Contains the logging setup and utilities.

## Accessing Swagger Documentation

The Swagger documentation for the API can be accessed at the following route:
```/swagger/```

This route provides a user-friendly interface to explore and test the API endpoints.

## Getting Started

To run the application, use the following command:

```
export CONFIG_PATH="./config/local.yaml"
go run cmd/server/main.go
```
Ensure you have the necessary dependencies installed by running
```
go mod tidy
```

Running with Docker

You can also run the application using Docker.

1. Build the Docker Image
```
docker build --build-arg CONFIG_PATH=/config/local.yaml -t ticket-processor .
```
2. Run the Container
```
docker run -e CONFIG_PATH=/config/local.yaml -p 8080:8080 ticket-processor
```
3. Access the API

Once the container is running, you can access the API at:
```
http://localhost:8080
```
And view the Swagger documentation at:
```
http://localhost:8080/swagger/
```


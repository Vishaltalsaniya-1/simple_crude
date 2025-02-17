# Simple CRUD Application

This is a simple CRUD (Create, Read, Update, Delete) application built using Go (Golang) with MongoDB as the database, RabbitMQ for message queuing, and Paota for distributed task processing.

## Features
- RESTful API with CRUD operations
- MongoDB for data storage
- RabbitMQ integration for async processing
- Paota worker pool for distributed task execution
- Dockerized setup

## Project Structure
```
simple_crude/
    ├── Dockerfile
    ├── config/           # Configuration files
    ├── consumer/         # RabbitMQ consumers
    ├── controller/       # API request handlers
    ├── db/              # Database connection setup
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── main.go          # Entry point
    ├── manager/         # Business logic layer
    ├── models/          # Data models
    ├── producer/        # RabbitMQ producers
    ├── request/         # Request DTOs
    ├── response/        # Response DTOs
```

## Prerequisites
- Go 1.21+
- MongoDB
- RabbitMQ
- Docker (optional for containerized setup)

## Installation & Setup

### Clone the Repository
```sh
git clone https://github.com/Vishaltalsaniya-1/simple_crude.git
cd simple_crude
```

### Environment Configuration
Create a `.env` file in the root directory and configure the necessary variables:
```env
MONGO_URI=mongodb://localhost:27017
RABBITMQ_URI=amqp://guest:guest@localhost:5672/
PORT=8082
```

### Run with Docker
```sh
docker-compose up --build
```

### Run Locally
1. Start MongoDB and RabbitMQ manually.
2. Run the Go application:
```sh
go run main.go
```

## API Endpoints

### Create Resource
```http
POST /api/resource
```
**Request Body:**
```json
{
  "auth_id": "65c9f1a9bfe9f8b2d3a4c5e6", 
  "name":"pointery",
  "description": "A stdandard student profile",
  "tag": ["golang", "backend", "developer"],
  "student": true
}

```

### Read Resource
```http
GET /api/students
```

### Update Resource
```http
PUT /api/studentse/{id}
```

### Delete Resource
```http
DELETE /api/students/{id}
```

## RabbitMQ & Paota Integration
- Messages are published to RabbitMQ for asynchronous processing.
- Paota workers process queued tasks efficiently.

## Contributing
Feel free to fork the project and contribute!


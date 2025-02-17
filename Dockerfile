# Use the official Golang image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules and dependencies files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 8081

# Run the application
CMD ["./main"]

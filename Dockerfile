# Stage 1: Build the Go application
FROM golang:1.23.4-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service cmd/main.go

# Stage 2: Create a minimal Docker image for the application
FROM alpine:latest

# Install necessary runtime libraries
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/notification-service .

# Expose the port that the app will run on
EXPOSE 8080

# Run the application
CMD ["./notification-service"]
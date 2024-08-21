# Use the official Golang image to build the application
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal base image to reduce the final image size
FROM debian:bookworm-slim

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/. .

# Debug: List files to verify the binary exists
RUN ls -l /app

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
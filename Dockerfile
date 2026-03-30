# Use the official Golang image as the base
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

# Use a minimal base image for the final stage
FROM alpine:3.21

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]

# Use the official Golang image as the base
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application - main.go is in ./cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

# Use a minimal base image for the final stage
FROM alpine:3.21

# ✅ CA CERTIFICATES (in FINAL stage - critical for HTTPS)
RUN apk --no-cache add ca-certificates

# ✅ NON-ROOT USER (in FINAL stage - security requirement)
RUN adduser -D -g '' appuser
USER appuser

WORKDIR /home/appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Expose the application port
EXPOSE 8080

ENTRYPOINT ["./myapp"]

# Command to run the application
CMD []
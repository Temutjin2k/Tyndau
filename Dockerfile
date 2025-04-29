FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/notification

# Create final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/notification-service .

# Copy templates
COPY --from=builder /app/templates ./templates

# Expose ports if needed
# EXPOSE 8080

# Run the application
CMD ["./notification-service"]
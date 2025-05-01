# Stage 1: builder
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY notification-service/go.mod notification-service/go.sum ./
RUN go mod download

# Copy the entire notification-service folder
COPY notification-service/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd/notification


# Stage 2: minimal image
FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/templates

# Copy the binary from builder to the correct location
COPY --from=builder /notification-service /app/notification-service

# Optional: only if you have templates
COPY notification-service/templates/ /app/templates/

# Give execute permissions to the binary
RUN chmod +x /app/notification-service

CMD ["/app/notification-service"]

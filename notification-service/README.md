# Notification Service Architecture

## 📋 Overview

Notification Service is a specialized microservice designed to handle user communication through email notifications. It subscribes to specific events via NATS messaging system, processes them according to business rules, and delivers personalized notifications using HTML templates.


## Architecture
This diagram shows the architecture of the Notification Service.

![Notification Service Architecture](./notification-service/xdiagrams/diagram.svg)

The diagrams demonstrates how different services interact with each other and with external systems.

![Notification Service Architecture](./notification-service/xdiagrams/diagram1.svg)

![Notification Service Architecture](./notification-service/xdiagrams/diagram2.svg)

![Notification Service Architecture](./notification-service/xdiagrams/diagram3.svg)

### Key Features

- **Event-Driven Architecture**: Leverages NATS for reliable, scalable event processing
- **Templated Notifications**: Uses HTML templates to create consistent, branded emails
- **Multiple Event Types**: Currently supports `user.registered` and `music.album_released` events
- **Clean Architecture**: Follows domain-driven design principles for maintainability
- **Containerized**: Fully dockerized for easy deployment and scaling

### Component Breakdown

- **Event Consumer**: Subscribes to NATS subjects and receives events
- **Event Processor**: Validates events and determines appropriate actions
- **Template Engine**: Renders HTML templates with dynamic data
- **Email Sender**: Handles SMTP communication for reliable delivery


## 🔄 Event Flow

1. External services publish events to NATS (e.g., `user.registered`, `music.album_released`)
2. Notification service consumes these events
3. Events are validated and processed
4. Appropriate HTML template is selected and populated with event data
5. Email is formatted and sent via SMTP
6. Delivery status is logged

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- Docker and Docker Compose
- NATS server
- SMTP server access

## Env example
#### NATS Configuration
NATS_URL=nats://nats:4222
NATS_CLUSTER_ID=notification-cluster
NATS_CLIENT_ID=notification-service

#### SMTP Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
SMTP_FROM=notifications@example.com

#### Service Configuration
LOG_LEVEL=info
TEMPLATE_DIR=./templates

### Running with Docker

```shellscript
# Build and start the service
docker-compose up -d

# View logs
docker-compose logs -f notification-service

# Stop the service
docker-compose down
```

### Project Structure

```plaintext
notification-service/
├── cmd/                    # Application entry points
├── config/                 # Configuration management
├── internal/               # Private application code
│   ├── adapter/            # External service integrations
│   ├── app/                # Application service
│   ├── model/              # Domain models
│   └── usecase/            # Business logic
├── pkg/                    # Public libraries
├── templates/              # HTML email templates
└── xdiagrams/              # Architecture diagrams
```

### Testing

```shellscript
# Run all tests
make test

# Run specific tests
go test ./internal/usecase/...
```

## 📊 Monitoring and Observability

The service includes structured logging via the custom logger package. For production deployments, consider integrating with:

- Prometheus for metrics
- Jaeger or Zipkin for distributed tracing
- ELK stack for log aggregation

## 🔍 Architecture Diagrams

The `xdiagrams/` and `xmermaids/` directories contain detailed architecture diagrams that visualize:

- Overall service architecture
- Data flow between components
- Event processing sequence
- Integration with external systems

## Notes
The project uses PostgreSQL for possible data storage, but the current version does not use the database.

To send emails, you need to configure SMTP parameters correctly.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Contact

For questions or support, please contact: [beibarys7ergaliev@gmail.com](mailto:beibarys7ergaliev@gmail.com)

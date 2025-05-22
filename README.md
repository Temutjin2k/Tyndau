# Tyndau
Tyndau is a modern microservice-based music streaming platform designed to deliver a seamless listening experience. It offers features such as user authentication, music browsing, playback, and notifications.

<img width="1439" alt="Снимок экрана 2025-05-22 в 11 08 36" src="https://github.com/user-attachments/assets/7215cc6a-fc62-44be-8569-1fc70d0adcab" />

## Architecture
Tyndau follows a microservice architecture. Each component is isolated and communicates via gRPC, promoting scalability and maintainability.

Main Services:

* API Gateway – Central point routing incoming HTTP/gRPC requests to appropriate services.
* Auth Service – Handles user authentication, registration, token generation, and validation.
* User Service – Manages user profiles and preferences.
* Music Service – Manages music metadata, storage, search, and streaming logic.
* Notification Service – Sends user notifications (e.g., registration confirmation).
* Frontend – Basic HTML frontend for interaction with the platform.
  
All services are containerized using Docker and orchestrated via Docker Compose.

## Technologies Used

| Layer         | Technologies                                |
|---------------|---------------------------------------------|
| Language      | Go (Golang)                                 |
| Communication | gRPC                                        |
| Containers    | Docker, Docker Compose                      |
| Frontend      | HTML, basic CSS                             |
| Automation    | Makefile for builds and service management  |

## Project Structure

```bash
Tyndau/
├── api-gateway/
├── auth_service/
├── user_service/
├── music-service/
├── notification-service/
├── frontend/
├── proto/               
├── docker-compose.yml
└── Makefile
```

## Getting Started

1. Clone the repository
   
```bash
git clone --branch final --single-branch https://github.com/Temutjin2k/Tyndau.git
cd Tyndau
```

2. Create a .env file in each service directory (if applicable) based on provided .env.example files. Set up environment variables such as ports, database URLs, secret keys, etc.

3. Run all services

```bash
docker-compose up --build
```

This will build and start all microservices and the frontend.

### Proto Contracts
https://github.com/Temutjin2k/TyndauProto

### Students (SE-2308 group)
Merey Ibraim, Temutjin Koszhanov, Beibars Yergali


# Taskflow

# TL;DR

A simple event-driven task manager built with Go, using a microservices architecture.
Features include user authentication with JWT, secure inter-service communication via mutual TLS, Redis-based async notifications, and service discovery with Consul.
Tech stack: Go, gRPC, Docker, Redis, Vault, Consul.

[![Go Version](https://img.shields.io/badge/Go-1.23.6-blue.svg)](https://golang.org/)

A portfolio project demonstrating a microservices architecture built entirely in Go, leveraging HashiCorp Vault for secrets management and Consul for service discovery.

## Overview & Motivation

This project was created as a learning exercise and portfolio piece to explore and implement concepts related to:

- **Microservices Architecture:** Building independent, deployable services.
- **Go (Golang):** Using Go for building the backend services.
- **gRPC:** Implementing efficient service communication.
- **HashiCorp Vault:** Securely managing sensitive data like private keys for JWT signing.
- **HashiCorp Consul:** Implementing service discovery and registration.
- **Event-Driven Architecture:** Using Redis Pub/Sub for asynchronous communication between services.
- **Secure Development Practices:** Focusing on secure handling of credentials and tokens.

The goal is to simulate a simplified task management system where users can register, authenticate, and manage tasks—with notifications triggered for task assignments. As I learn new concepts, this project will be updated.

## Architecture

The system currently consists of three core microservices:

1. **User Service:**
    - Manages user registration and authentication.
    - Generates JWTs upon successful login/registration.
    - Signs JWTs using a private key securely retrieved from HashiCorp Vault's KV store.
    - Provides user details (like email) to other services.
2. **Task Service:**
    - Manages the creation and retrieval of tasks.
    - Allows creating tasks assigned to specific users or unassigned tasks.
    - Lists tasks assigned to a user or retrieves all tasks.
    - Publishes an event to a Redis channel when a new task is assigned to a user.
3. **Notifier Service:**
    - Subscribes to the task assignment event channel on Redis.
    - Upon receiving an event, retrieves the relevant user's email from the User service via gRPC (using Consul for discovery).
    - Sends an email notification to the user about their newly assigned task.

**Communication:**

- Services communicate with each other using **gRPC**.
- Services register themselves with **Consul** upon startup.
- The Notifier service currently discovers User service instances via Consul and selects one **randomly** for communication. (Note: This is a simplification for learning purposes).
- The Task service publishes events to **Redis**, and the Notifier service subscribes to these events.

## Technology Stack

- **Language:** Go (Golang)
- **Frameworks/Libraries:**
  - gRPC
  - HashiCorp Vault API
  - HashiCorp Consul API
  - Redis Client
  - JWT Library
  - Go DotEnv
- **Infrastructure:**
  - HashiCorp Vault (KV Secret Engine v2)
  - HashiCorp Consul
  - Redis
- **Other:** Docker & Docker Compose (Recommended for local setup)

## Setup & Installation

**Prerequisites:**

- Go (Version 1.23.6 or later)
- Docker & Docker Compose (Recommended for running Vault, Consul, Redis)
- Vault CLI (Optional, for manual interaction)
- Consul CLI (Optional, for manual interaction)
- Gmail App Password

**Steps:**

1. **Clone the repository:**

    ```bash
    git clone https://github.com/[Your GitHub Username]/[Your Project Name].git
    cd [Your Project Name]
    ```

2. **Configure Environment Variables:**

    - Copy the example environment file:

      ```bash
      # Each service has its own env file.
      # There is a global file in ./config
      cp .env.example .env
      ```

    - The environment files may include values such as:
      - Vault Address (`VAULT_ADDR`)
      - Approle Client ID (`APPROLE_ROLE_ID`)
      - Approle Secret (`APPROLE_SECRET_ID`)
      - Vault Secret Path for JWT Key (`VAULTKEY_PATH`) - e.g., `data/jwt/auth`
      - Vault Secret Key Name for JWT Key (`VAULTKEY_NAME`) - e.g., `private_key`
      - Redis Address (`REDIS_NOTIFIER_ADDR`)
      - Email to send notification from (`GMAIL_SOURCE`)
      - Gmail App Password (`GMAIL_APP_PASSWORD`)

3. **Start Infrastructure (Vault, Consul, Redis):**
    - Use docker-compose to setup HashiCorp Vault and Redis

      ```bash
      cd scripts
      docker-compose up -d
      ```

    - Run Consul

      ```bash
      docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
      ```

4. **Configure Vault:**

    - Ensure the Vault KV v2 engine is enabled at the path specified in `.env` (usually `secret/` by default in dev mode).
    - Store your private key in Vault at the path specified (`VAULT_KEY_PATH`) with the key name (`VAULT_KEY_NAME`). _You will need to generate a private/public key pair first (e.g., using `openssl`)._

    ```bash
    # Example using Vault CLI (ensure VAULT_ADDR and VAULT_TOKEN are set)
    # vault kv put <VAULT_KV_MOUNT>/<VAULT_JWT_KEY_PATH> <VAULT_JWT_KEY_NAME>=@"path/to/your/local/private_key.pem"
    vault kv put secret/data/jwt/auth private_key=@/path/to/your/private_key.pem
    ```

    - _Note: Ensure your public key is available separately if needed later for the gateway._

5. **Build and Run Services:**

    - Navigate to each service directory and run:

      ```bash
      # Example for User Service
      cd user/cmd/server/
      go run main.go # optional flag -port

      # Open new terminals for the other services
      cd task/cmd/server/
      go run main.go # optional flag -port

      cd notifier/cmd/server/
      go run main.go # optional flag -port
      ```

## Usage

Since the services communicate via gRPC and there is no API Gateway yet, you'll need a gRPC client (like `grpcurl`, Evans, or Postman's gRPC feature) to interact with them directly. (Postman recommended)
If using Postman, create a new gRPC collection and upload the provided `.proto` files: `./api/task/v1/task.proto` and `./api/user/v1/user.proto`.
Click `Use Example Message`, fill in the request body, and click `Send`.

## Future Enhancements / To-Do

This project serves as a foundation. Planned future improvements include:

- **Vault Transit Engine:** Replace Vault KV storage of the private key with the Transit engine for signing operations, preventing the key from being exposed in memory.
- **Task Assignment Endpoint:** Implement a dedicated /assignTask endpoint in the Task service.
- **Improved Service Discovery:** Implement more robust service instance selection from Consul (e.g., load balancing strategies like round-robin instead of random).
- **API Gateway:** Introduce a gateway service that exposes RESTful or GraphQL endpoints to external clients (e.g., a frontend) and communicates with backend services via gRPC.
- **mTLS Implementation:** Secure inter-service gRPC communication using mutual TLS (mTLS), potentially using Vault as the Certificate Authority (CA).
- **Gateway JWT Verification:** The API Gateway should retrieve the public key from Vault to verify JWTs received from clients (which were originally signed by the User service).
- **Testing:** Add comprehensive unit and integration tests.
- **Containerization:** Improve Dockerfiles for production-like builds and optimize Docker Compose setup.
- **Observability:** Add logging, metrics, and tracing.
- Persistent Storage: Replace in-memory storage with a proper database (e.g., MySQL or PostgreSQL). Current services are designed around interfaces to make this transition seamless.

## Security Considerations

- **JWT Key Storage:** Currently uses Vault KV. While more secure than storing keys in code or configuration files, the key is read into the User service's memory. Migrating to Vault Transit Engine is recommended for higher security.
- **Service Communication:** Currently relies on plaintext gRPC. Implementing mTLS is crucial for securing inter-service communication in a real-world scenario.
- **Authentication/Authorization:** Basic JWT authentication is implemented. More robust authorization logic (e.g., ensuring only the assigned user can modify their tasks) should be added.
- **Secret Management:** Ensure Vault tokens and other sensitive configurations are managed securely (e.g., not hardcoded, using appropriate Vault policies).

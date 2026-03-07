# go-infra

`go-infra` is a Go-based infrastructure utility service designed to provide common shared functionalities for distributed systems and microservices. It centralizes messaging (SMS/Email), configuration distribution, health monitoring, and internationalization.

This service is designed to be a lightweight, high-performance "sidekick" or central utility that other services in your ecosystem can call upon.

## Features

- **Messaging Gateway**: 
  - Asynchronous SMS and Email delivery via internal task queues.
  - Template-based email rendering (embedded HTML templates).
  - Pluggable HTTP-based providers (SMS/Email gateways).
- **Configuration Management**:
  - Serves static configuration files over HTTP to other services.
  - Multi-source configuration loading (Command-line flags, Environment variables, JSON files, and Remote URLs).
  - Support for environment variable expansion within configuration files.
- **Internationalization (i18n)**:
  - Centralized translation management.
  - Dynamic language switching and placeholder replacement.
- **System Monitoring**:
  - Prometheus metrics integration.
  - Dedicated system listener for metrics and management (optional separate port).
  - Health checks for database connectivity and connection pooling stats.
- **Reliability**:
  - Graceful shutdown support for clean connection termination.
  - Built-in task queue with worker limits and panic recovery.
  - Robust HTTP transport tuning (Idle connections, timeouts, etc.).

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: [Echo v4](https://echo.labstack.com/)
- **Database ORM**: [GORM](https://gorm.io/) (PostgreSQL)
- **Monitoring**: [Prometheus](https://prometheus.io/)
- **Logging**: [Structured Logging (slog)](https://pkg.go.dev/log/slog)

## Getting Started

### Prerequisites

- Go installed on your machine.
- A running PostgreSQL instance.
- (Optional) Docker for containerized deployment.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/go-infra.git
   cd go-infra
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o go-infra cmd/go-infra/main.go
   ```

### Running the Service

You can run the service using command-line flags or environment variables:

```bash
./go-infra -listen 0.0.0.0:30780 -env production
```

## Configuration

The service uses a hierarchical configuration approach (Flags > Env Vars > JSON Config).

### Environment Variables

Prefix your environment variables with `APP_`. Examples:
- `APP_ENV`: Environment mode (`development`, `production`, etc.)
- `APP_DB_HOST`: Database host address.
- `APP_DB_PASSWORD`: Database password.
- `APP_CONFIG`: Path to the directory containing `config.{env}.json`.
- `APP_SYS_API_KEY`: Security key for accessing sensitive `/sys` endpoints.

### File-based Configuration

The service looks for a `config.{env}.json` in the paths specified by the `APP_CONFIG` environment variable.

## API Endpoints

### Internal Messaging
- `POST /sys/api/messenger/sms-text`: Send a plain text SMS.
- `POST /sys/api/messenger/sms-passcode`: Send a localized 2FA passcode via SMS.
- `POST /sys/api/messenger/email-html`: Send a raw HTML email.
- `POST /sys/api/messenger/email-passcode`: Send a templated 2FA passcode via Email.

### Infrastructure & Health
- `GET /health`: Basic service liveness check.
- `GET /infra/api/ping`: Basic connectivity test.
- `GET /sys/api/metrics`: Prometheus metrics (Requires `APP_SYS_API_KEY`).
- `GET /sys/api/configs/*`: Serve static configuration files from the configured directory.

## Project Structure

```text
├── cmd/
│   └── go-infra/          # Main entry point
├── internal/
│   ├── cmd/               # Command execution logic and graceful shutdown
│   ├── config/            # Configuration loading and validation
│   ├── controller/        # Web handlers (Health, Messaging, etc.)
│   ├── i18n/              # Translation and localization logic
│   ├── repository/        # Database abstraction (GORM)
│   ├── router/            # Route definitions and middleware setup
│   ├── service/           # Business logic and messaging gateways
│   └── util/              # Shared utilities (Task queues, HTTP, logging)
└── test/
    └── e2e/               # End-to-end integration tests
```

## Deployment

### Docker
A sample service definition for Docker Compose:

```yaml
services:
  go-infra:
    image: alpine:3.20
    container_name: go-infra
    command: ./go-infra
    environment:
      - APP_ENV=production
      - APP_DB_HOST=postgres-db
      - APP_DB_NAME=infra_db
      - APP_DB_PASSWORD_FILE=/run/secrets/db_password
    volumes:
      - ./configs:/app/configs:ro
      - /etc/ssl/certs:/etc/ssl/certs:ro
    working_dir: /app
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
# go-infra

`go-infra` is a microservice built with Go and the [Echo](https://echo.labstack.com/) framework. It manages infrastructure-level settings, general application metadata, and common cross-cutting concerns for the architecture.

## Features

- **Database Management:** Uses GORM with PostgreSQL.
- **HTTP Server:** Built on the Echo framework.
- **Metrics:** Exposes Prometheus metrics for monitoring.

## Prerequisites

- Go 1.26+
- Python 3.x

## Build and Run

```sh
# Run tests
python Makefile.py test

# Run linter
python Makefile.py lint

# Build binary for Linux
python Makefile.py linux
```

## Architecture Context

This service provides foundational data and configurations required by other microservices in the deployment. Refer to the Docker Compose examples in `projects/ecom-shop` to see how it integrates with the rest of the application.

# Content Lambda Service

A Go-based AWS Lambda service for content management using DynamoDB.

## Project Structure

```
.
├── cmd/                # Main application entry points
│   └── lambda/        # Lambda handler
├── internal/          # Private application code
├── api/              # API definitions and interfaces
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Local development environment
├── go.mod            # Go module definition
└── Makefile          # Common development commands
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- AWS CLI (for local development)

## Local Development

1. Start the local development environment:
```bash
make docker-run
```

2. Initialize the DynamoDB table:
```bash
make init-db
```

3. Build and run the application locally:
```bash
make build
make run
```

## Docker Commands

- Build the Docker image:
```bash
make docker-build
```

- Start the development environment:
```bash
make docker-run
```

- Stop the development environment:
```bash
make docker-down
```

## Testing

Run the test suite:
```bash
make test
```

## Cleanup

Remove build artifacts:
```bash
make clean
``` 
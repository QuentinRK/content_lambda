FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/lambda

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Install AWS CLI and other necessary tools
RUN apk add --no-cache \
    aws-cli \
    curl \
    jq

# Copy the binary from builder
COPY --from=builder /app/main .

# Set the entrypoint
ENTRYPOINT ["./main"] 
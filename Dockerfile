# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Add a build argument for debug mode
ARG DEBUG=false

# Set build tags based on the DEBUG flag
RUN if [ "$DEBUG" = "true" ]; then \
      go mod download && go build -o server -tags debug ./cmd/server/main.go; \
    else \
      go mod download && go build -o server ./cmd/server/main.go; \
    fi

# Stage 2: Create the production image
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/server .

# Environment variables
ENV LISTEN_ADDR=:8080
ENV STORAGE_PATH=/app/storage

# Create the storage directory
RUN mkdir -p /app/storage

EXPOSE 8080

# Command to run the server
CMD ["/app/server"]

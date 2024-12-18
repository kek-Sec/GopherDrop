# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Add build arguments for debug mode and versioning
ARG DEBUG=false
ARG GIT_COMMIT_SHA
ARG GIT_VERSION

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

# Add OCI Image Spec labels
LABEL org.opencontainers.image.title="GopherDrop Backend" \
      org.opencontainers.image.description="Backend for GopherDrop, a secure one-time secret sharing service" \
      org.opencontainers.image.source="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.revision="${GIT_COMMIT_SHA}" \
      org.opencontainers.image.version="${GIT_VERSION}" \
      org.opencontainers.image.url="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.documentation="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.licenses="MIT"

# Environment variables
ENV LISTEN_ADDR=:8080
ENV STORAGE_PATH=/app/storage

# Create the storage directory
RUN mkdir -p /app/storage

EXPOSE 8080

# Command to run the server
CMD ["/app/server"]

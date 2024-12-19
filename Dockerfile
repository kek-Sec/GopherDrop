# Stage 1: Build the Go Backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app
COPY . .

# Generate Go version tag based on the current date in ddMMyyyy format
ARG GO_VERSION_TAG=$(date +"%d%m%Y")
ENV GO_VERSION_TAG=${GO_VERSION_TAG}

# Add build arguments for debug mode and versioning
ARG DEBUG=false
ARG GIN_MODE=release
ARG GIT_COMMIT_SHA
ARG GIT_VERSION
ENV GIN_MODE=${GIN_MODE}

# Set build tags based on the DEBUG flag and include versioning information
RUN if [ "$DEBUG" = "true" ]; then \
      go mod download && go build -o server -tags debug -ldflags="-X main.version=DEBUG" ./cmd/server/main.go; \
    else \
      go mod download && go build -o server -ldflags="-X main.version=${GO_VERSION_TAG}-${GIT_VERSION}-${GIT_COMMIT_SHA}" ./cmd/server/main.go; \
    fi

# Stage 2: Build the Vue.js Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app
COPY ui/package.json ui/package-lock.json ./ 
RUN npm install --legacy-peer-deps

# Add build argument for the API URL and versioning
ARG VITE_API_URL="/api"
ARG GIT_VERSION
ARG GIT_COMMIT_SHA
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_APP_VERSION=${GIT_VERSION}-${GIT_COMMIT_SHA}

COPY ui ./ 
RUN npm run build

# Stage 3: Combine Backend and Frontend into a Single Image
FROM nginx:alpine

# Add OCI Image Spec labels
ARG GIT_COMMIT_SHA
ARG GIT_VERSION

LABEL org.opencontainers.image.title="GopherDrop" \
      org.opencontainers.image.description="GopherDrop - Secure one-time secret sharing service" \
      org.opencontainers.image.source="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.revision="${GIT_COMMIT_SHA}" \
      org.opencontainers.image.version="${GIT_VERSION}" \
      org.opencontainers.image.url="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.documentation="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.licenses="MIT"

# Copy the Go server binary
COPY --from=backend-builder /app/server /app/server

# Copy the frontend static files
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create the storage directory for the backend
RUN mkdir -p /app/storage

# Expose the ports for Nginx and the Go server
EXPOSE 80 8080

# Run both the Go server and Nginx using a simple script
CMD ["/bin/sh", "-c", "/app/server & nginx -g 'daemon off;'"]

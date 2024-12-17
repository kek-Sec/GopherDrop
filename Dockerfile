FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o server ./cmd/server/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
ENV LISTEN_ADDR=:8080
ENV STORAGE_PATH=/app/storage
RUN mkdir -p /app/storage
EXPOSE 8080
CMD ["/app/server"]

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY . .

RUN go build -o bin/proxy ./cmd/proxy

# Execute stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/proxy ./

EXPOSE 8080

CMD ["./proxy"]
# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the furnace worker
RUN CGO_ENABLED=0 GOOS=linux go build -o furnace-worker ./cmd/workerserver/furnace

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/furnace-worker .

CMD ["./furnace-worker"]

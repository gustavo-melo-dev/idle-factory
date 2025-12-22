# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the drill worker
RUN CGO_ENABLED=0 GOOS=linux go build -o drill-worker ./cmd/workerserver/drill

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/drill-worker .

CMD ["./drill-worker"]

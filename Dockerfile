# ---------- Step 1: Build ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app

# ---------- Step 2: Run ----------
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Run
CMD ["./app"]
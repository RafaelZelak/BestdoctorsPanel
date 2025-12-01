# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Tidy dependencies after Go version change
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bestdoctors_service cmd/main.go

# Final stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bestdoctors_service .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 9002

# Run the application
CMD ["./bestdoctors_service"]

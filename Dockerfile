# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata protobuf protobuf-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install protobuf tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Generate protobuf files
RUN mkdir -p gen/proto
RUN protoc -I proto --go_out=gen/proto --go_opt=paths=source_relative --go-grpc_out=gen/proto --go-grpc_opt=paths=source_relative --grpc-gateway_out=gen/proto --grpc-gateway_opt=paths=source_relative proto/captcha/captcha.proto
RUN protoc -I proto --go_out=gen/proto --go_opt=paths=source_relative --go-grpc_out=gen/proto --go-grpc_opt=paths=source_relative proto/balancer/balancer.proto

# Build captcha-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o captcha-service ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl tzdata

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/captcha-service .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/backgrounds ./backgrounds

# Create logs directory
RUN mkdir -p logs

# Expose port
EXPOSE 8083

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8083/health || exit 1

# Run captcha-service
CMD ["./captcha-service"]
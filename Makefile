.PHONY: proto build run test test-unit test-integration test-captcha test-performance test-bot-protection clean docker-build docker-run

proto:
	@echo "Generating protobuf files..."
	@mkdir -p gen/proto/balancer gen/proto/captcha
	@protoc --go_out=gen/proto/balancer --go_opt=paths=source_relative \
		--go-grpc_out=gen/proto/balancer --go-grpc_opt=paths=source_relative \
		proto/balancer/balancer.proto
	@protoc --go_out=gen/proto/captcha --go_opt=paths=source_relative \
		--go-grpc_out=gen/proto/captcha --go-grpc_opt=paths=source_relative \
		proto/captcha/captcha.proto
	@echo "Protobuf files generated successfully"

build: proto
	@echo "Building captcha service..."
	@go build -o bin/captcha-service ./cmd/server
	@echo "Build completed"

run: build
	@echo "Starting captcha service..."
	@./bin/captcha-service

test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	@go test -v ./internal/... ./pkg/...

test-integration: test-captcha test-performance test-bot-protection

test-captcha:
	@echo "Running captcha integration tests..."
	@go test -v -run TestCaptchaServiceIntegration ./tests/integration/

test-performance:
	@echo "Running performance tests..."
	@go test -v -run TestCaptchaServicePerformance ./tests/integration/

test-bot-protection:
	@echo "Running bot protection tests..."
	@go test -v -run TestBotProtection ./tests/integration/

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf gen/
	@go clean
	@docker system prune -f

docker-build:
	@echo "Building Docker image..."
	@docker build -t captcha-service .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 38000:38000 captcha-service

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

dev:
	@echo "Starting development server..."
	@go run ./cmd/server

help:
	@echo "Available commands:"
	@echo "  proto              - Generate protobuf files"
	@echo "  build              - Build the application"
	@echo "  run                - Build and run the application"
	@echo "  test               - Run all tests (unit + integration)"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run all integration tests"
	@echo "  test-captcha       - Run captcha functionality tests"
	@echo "  test-performance   - Run performance tests"
	@echo "  test-bot-protection - Run bot protection tests"
	@echo "  clean              - Clean build artifacts and Docker"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Build and run Docker container"
	@echo "  install-deps       - Install dependencies"
	@echo "  dev                - Run in development mode"
	@echo "  help               - Show this help message"
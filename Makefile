.PHONY: build run clean test install

# Binary name
BINARY_NAME=sshresume
OUTPUT_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/sshresume

# Run the application
run: build
	@echo "Starting SSH server..."
	./$(OUTPUT_DIR)/$(BINARY_NAME)

# Run directly without building (uses go run)
dev:
	@echo "Running in development mode..."
	go run ./cmd/sshresume --host localhost --port 23234

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Display help
help:
	@echo "Available commands:"
	@echo "  make build       - Build the binary to $(OUTPUT_DIR)/$(BINARY_NAME)"
	@echo "  make run         - Build and run the server (localhost:23234)"
	@echo "  make run-custom  - Run with custom HOST and PORT (e.g., make run-custom HOST=0.0.0.0 PORT=2222)"
	@echo "  make dev         - Run directly with 'go run' (faster for development)"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make install     - Install/update dependencies"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make help        - Show this help message"

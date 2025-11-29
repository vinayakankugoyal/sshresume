.PHONY: build run clean test install build-linux-arm64

# Binary name
BINARY_NAME=sshresume
OUTPUT_DIR=bin

# Build the application for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/sshresume

# Build for Linux AMD64
build-linux-amd64:
	@echo "Building $(BINARY_NAME) for Linux AMD64..."
	@mkdir -p $(OUTPUT_DIR)/linux/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/linux/amd64/$(BINARY_NAME) ./cmd/sshresume

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
	@rm -rf .ssh

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

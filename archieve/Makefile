# Robot Arm Controller Makefile

.PHONY: help build run test clean install-deps fmt vet run-rpi run-coral deploy deploy-coral deploy-coral-interactive interactive

# Default target
help:
	@echo "Robot Arm Controller - Available targets:"
	@echo "  build        Build the arm controller binary"
	@echo "  run          Run the arm controller demo (mock board)"
	@echo "  interactive  Run in interactive mode for testing"
	@echo "  test-display-coral  Test display functionality with Coral board"
	@echo "  run-rpi      Run with Raspberry Pi configuration"
	@echo "  run-coral    Run with Google Coral configuration"
	@echo "  run-linux    Run with generic Linux GPIO configuration"
	@echo "  test         Run all tests"
	@echo "  clean        Clean build artifacts"
	@echo "  install-deps Install Go dependencies"
	@echo "  fmt          Format Go code"
	@echo "  vet          Run go vet"
	@echo "  all          Run fmt, vet, test, and build"
	@echo "  deploy       Deploy and run on remote device (usage: make deploy SSH_HOST=<ssh_config_name> [BOARD=<board_type>] [CONFIG=<config_type>])"
	@echo "  deploy-coral Deploy to Coral board (usage: make deploy-coral SSH_HOST=<ssh_config_name>)"
	@echo "  deploy-coral-interactive  Deploy and run interactive mode on Coral board"

# Build the main binary
build:
	mkdir -p ./bin
	go build -o ./bin/arm-controller ./cmd/arm

# Run the demo application with different board configurations
run:
	go run ./cmd/arm -board=mock -verbose=true

run-rpi:
	go run ./cmd/arm -board=raspberry-pi -config=raspberry-pi -verbose=true

run-coral:
	go run ./cmd/arm -board=coral -config=coral -verbose=true

run-linux:
	go run ./cmd/arm -board=linux-gpio -config=generic -verbose=true

# Cross-compile for Raspberry Pi (ARM64)
build-rpi:
	mkdir -p ./bin
	GOOS=linux GOARCH=arm64 go build -o ./bin/arm-controller-rpi ./cmd/arm

# Cross-compile for Raspberry Pi (ARM32)
build-rpi32:
	mkdir -p ./bin
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/arm-controller-rpi32 ./cmd/arm

# Run tests
test:
	go test -v ./pkg/... ./internal/... ./cmd/...

# Clean build artifacts
clean:
	rm -f ./bin/arm-controller ./bin/arm-controller-rpi ./bin/arm-controller-rpi32 arm-controller
	go clean

# Install dependencies
install-deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks and build
all: fmt vet test build
	@echo "✅ All checks passed and binary built successfully!"

# Install Go (macOS specific)
install-go:
	@echo "Installing Go..."
	@if command -v brew >/dev/null 2>&1; then \
		brew install go; \
	else \
		@echo "Homebrew not found. Please install Go manually from https://golang.org/dl/"; \
	fi

# Development setup
setup: install-go install-deps
	@echo "🚀 Development environment setup complete!"

# Run examples
run-display-demo:
	cd examples && go run robot-display-demo.go

run-6dof-demo:
	cd examples && go run 6dof-pan-tilt-config.go

run-coral-example:
	cd examples && go run coral-config.go

run-interactive:
	cd examples && go run interactive-display-demo.go

run-advanced:
	cd examples && go run advanced-display-demo.go

run-interactive-demo:
	cd examples && go run interactive-display-demo.go

run-advanced-demo:
	cd examples && go run advanced-display-demo.go

# Deploy to remote device via SSH
deploy:
	@if [ -z "$(SSH_HOST)" ]; then \
		echo "❌ Error: SSH_HOST is required. Usage: make deploy SSH_HOST=<ssh_config_name>"; \
		echo "Example: make deploy SSH_HOST=raspberry-pi BOARD=raspberry-pi CONFIG=raspberry-pi"; \
		exit 1; \
	fi
	@echo "🚀 Deploying to $(SSH_HOST)..."
	@echo "📦 Building binary for Linux ARM64..."
	mkdir -p ./bin
	GOOS=linux GOARCH=arm64 go build -o ./bin/arm-controller-deploy ./cmd/arm
	@echo "📤 Transferring binary to $(SSH_HOST)..."
	scp ./bin/arm-controller-deploy $(SSH_HOST):~/arm-controller
	@echo "🔧 Making binary executable..."
	ssh $(SSH_HOST) "chmod +x ~/arm-controller"
	@echo "▶️  Running arm controller on $(SSH_HOST)..."
	@if [ -n "$(BOARD)" ] && [ -n "$(CONFIG)" ]; then \
		echo "Running with board=$(BOARD) and config=$(CONFIG)"; \
		ssh $(SSH_HOST) "cd ~ && ./arm-controller -board=$(BOARD) -config=$(CONFIG) -verbose=true"; \
	elif [ -n "$(BOARD)" ]; then \
		echo "Running with board=$(BOARD)"; \
		ssh $(SSH_HOST) "cd ~ && ./arm-controller -board=$(BOARD) -verbose=true"; \
	else \
		echo "Running with default configuration"; \
		ssh $(SSH_HOST) "cd ~ && ./arm-controller -verbose=true"; \
	fi
	@echo "✅ Deployment and execution completed!"

# Quick deploy to Coral robot arm (convenience target)
deploy-coral:
	@if [ -z "$(SSH_HOST)" ]; then \
		echo "❌ Error: SSH_HOST is required. Usage: make deploy-coral SSH_HOST=<ssh_config_name>"; \
		echo "Example: make deploy-coral SSH_HOST=robot-arm"; \
		echo "This is equivalent to: make deploy SSH_HOST=robot-arm BOARD=coral CONFIG=coral"; \
		exit 1; \
	fi
	@echo "🚀 Quick deploying to Coral board at $(SSH_HOST)..."
	$(MAKE) deploy SSH_HOST=$(SSH_HOST) BOARD=coral CONFIG=coral

# Deploy with interactive mode to Coral robot arm
deploy-coral-interactive:
	@if [ -z "$(SSH_HOST)" ]; then \
		echo "❌ Error: SSH_HOST is required. Usage: make deploy-coral-interactive SSH_HOST=<ssh_config_name>"; \
		echo "Example: make deploy-coral-interactive SSH_HOST=robot-arm"; \
		exit 1; \
	fi
	@echo "🚀 Deploying interactive mode to Coral board at $(SSH_HOST)..."
	@echo "📦 Building binary for Linux ARM64..."
	mkdir -p ./bin
	GOOS=linux GOARCH=arm64 go build -o ./bin/arm-controller-deploy ./cmd/arm
	@echo "📤 Transferring binary to $(SSH_HOST)..."
	scp ./bin/arm-controller-deploy $(SSH_HOST):~/arm-controller
	@echo "🔧 Making binary executable..."
	ssh $(SSH_HOST) "chmod +x ~/arm-controller"
	@echo "🎮 Starting interactive mode on $(SSH_HOST)..."
	ssh -t $(SSH_HOST) "cd ~ && ./arm-controller -board=coral -config=coral -interactive=true -verbose=true"

# Run in interactive mode for testing and development
interactive:
	@echo "🎮 Starting interactive mode (use 'exit' or 'quit' to exit properly)"
	go run ./cmd/arm -board=mock -interactive=true -verbose=true

interactive-rpi:
	@echo "🎮 Starting interactive mode for Raspberry Pi"
	go run ./cmd/arm -board=raspberry-pi -config=raspberry-pi -interactive=true -verbose=true

interactive-coral:
	@echo "🎮 Starting interactive mode for Coral Dev Board"
	go run ./cmd/arm -board=coral -config=coral -interactive=true -verbose=true

# Quick test of interactive mode functionality
test-interactive:
	@echo "🧪 Testing interactive mode (non-interactive demo)"
	go run ./cmd/arm -board=mock -interactive=false -verbose=false

# Test display functionality with Coral board
test-display-coral:
	@echo "📺 Testing display with Coral board configuration"
	@echo "Commands: display Hello, display status, display test, exit"
	go run ./cmd/arm -board=coral -config=coral -interactive=true -verbose=true

# Robot Arm Controller Makefile

.PHONY: help build run test clean install-deps fmt vet run-rpi run-coral

# Default target
help:
	@echo "Robot Arm Controller - Available targets:"
	@echo "  build        Build the arm controller binary"
	@echo "  run          Run the arm controller demo (mock board)"
	@echo "  run-rpi      Run with Raspberry Pi configuration"
	@echo "  run-coral    Run with Google Coral configuration"
	@echo "  run-linux    Run with generic Linux GPIO configuration"
	@echo "  test         Run all tests"
	@echo "  clean        Clean build artifacts"
	@echo "  install-deps Install Go dependencies"
	@echo "  fmt          Format Go code"
	@echo "  vet          Run go vet"
	@echo "  all          Run fmt, vet, test, and build"

# Build the main binary
build:
	go build -o arm-controller ./cmd/arm

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
	GOOS=linux GOARCH=arm64 go build -o arm-controller-rpi ./cmd/arm

# Cross-compile for Raspberry Pi (ARM32)
build-rpi32:
	GOOS=linux GOARCH=arm GOARM=7 go build -o arm-controller-rpi32 ./cmd/arm

# Run tests
test:
	go test -v ./pkg/... ./internal/... ./cmd/...

# Clean build artifacts
clean:
	rm -f arm-controller
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

.PHONY: build install clean deps playwright-install uninstall

# Binary name
BINARY_NAME=bibapa

# Install location
INSTALL_PATH=/usr/local/bin

# Build the binary
build:
	go build -o $(BINARY_NAME) ./cmd/bibapa

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install Playwright browsers (one-time setup)
playwright-install:
	go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps chromium

# Install the binary to system path
install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_PATH)"

# Uninstall the binary from system path
uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME) from $(INSTALL_PATH)"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Full setup: deps + playwright + build + install
setup: deps playwright-install install
	@echo "Setup complete!"

# Development build with race detector
build-dev:
	go build -race -o $(BINARY_NAME) ./cmd/bibapa

# Run tests
test:
	go test -v ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  install            - Build and install to $(INSTALL_PATH)"
	@echo "  uninstall          - Remove from $(INSTALL_PATH)"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  playwright-install - Install Playwright browsers (one-time)"
	@echo "  setup              - Full setup (deps + playwright + install)"
	@echo "  clean              - Remove build artifacts"
	@echo "  build-dev          - Build with race detector"
	@echo "  test               - Run tests"
	@echo "  help               - Show this help"


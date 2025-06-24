.PHONY: build install clean test run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=xsh
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Install the binary to GOBIN
install: build
	$(GOCMD) install .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run the application
run: build
	./$(BINARY_NAME)

# Cross compilation for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v .

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Check for go.mod and go.sum consistency
mod-verify:
	$(GOMOD) verify

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run go vet
vet:
	$(GOCMD) vet ./...

# Run all checks
check: fmt vet test

# Development setup
dev-setup:
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint code
lint:
	golangci-lint run

# Complete build and check
all: clean deps fmt vet test build

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Install the binary to GOBIN"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  run        - Build and run the application"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  check      - Run fmt, vet, and test"
	@echo "  lint       - Run golangci-lint"
	@echo "  all        - Clean, deps, fmt, vet, test, and build"
	@echo "  help       - Show this help message" 
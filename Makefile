.PHONY: build test lint fmt clean help

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -w .

# Clean build artifacts
clean:
	go clean ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the project"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linter"
	@echo "  fmt    - Format code with gofmt"
	@echo "  clean  - Clean build artifacts"
	@echo "  help   - Show this help message"
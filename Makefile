.PHONY: build test lint clean help

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	go clean ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the project"
	@echo "  test   - Run tests"
	@echo "  lint   - Run linter"
	@echo "  clean  - Clean build artifacts"
	@echo "  help   - Show this help message"
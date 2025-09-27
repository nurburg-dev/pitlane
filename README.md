# Pitlane

[![Test and Coverage](https://github.com/nurburg-dev/pitlane/actions/workflows/test.yml/badge.svg)](https://github.com/nurburg-dev/pitlane/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nurburg-dev/pitlane/branch/main/graph/badge.svg)](https://codecov.io/gh/nurburg-dev/pitlane)

A workflow orchestration engine with PostgreSQL backend.

## Features

- Workflow and activity execution
- PostgreSQL database with optimized indexes
- Transaction-based repositories
- Comprehensive test coverage
- Readable ID generation

## Development

### Prerequisites
- Go 1.24+
- PostgreSQL 15+
- Docker (for tests)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Database

The project uses PostgreSQL with embedded schema migrations and testcontainers for testing.
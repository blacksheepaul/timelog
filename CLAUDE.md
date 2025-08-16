# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based time logging application that allows users to track time entries with tags and remarks. It's a single-user local system with no authentication, built using a standard MVC architecture with Gin framework, SQLite database, and GORM ORM.

## Common Development Commands

### Building
```bash
# Build for current platform
make build

# Build for Linux ARM64 (for Docker)
make build-linux

# Build and create Docker image
make docker env=prod
make docker env=test
```

### Running
```bash
# Local development
./main

# Docker containers
make run env=prod    # Runs on port 8080
make run env=test    # Runs on port 18080

# Manual Docker run
docker run --rm -e ENV=prod -p 8080:8080 timelog-app
docker run --rm -e ENV=test -p 18080:8080 timelog-app-test
```

### Database Migrations
```bash
# Create new migration
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ init_xxx_table

# Apply migrations
migrate -database "sqlite3://dev.db" --path model/migrations/ up
```

### Testing
```bash
# Run integration tests
go test ./test/...

# The test setup uses config-test.yml and can flush the database if configured
```

## Architecture

### Layer Structure
- **main.go**: Application entry point with graceful shutdown
- **router/**: HTTP routing and handlers
  - `router.go`: Main router setup with middleware
  - `timelog.go`: TimeLog API endpoints (CRUD operations)
  - `middleware/`: CORS and auth middleware
- **service/**: Business logic layer
  - `timelog.go`: TimeLog business operations
- **model/**: Data access layer with GORM
  - `timelog.go`: TimeLog model and CRUD functions
  - `model.go`: Database connection and DAO pattern
  - `migrations/`: SQL migration files
- **core/**: Shared utilities
  - `config/`: Configuration management with cleanenv
  - `logger/`: Zap logger setup with file rotation

### Key Components
- **Database**: SQLite with CGO-free driver (`ncruces/go-sqlite3`)
- **Web Framework**: Gin with JSON logging formatter
- **ORM**: GORM with soft deletes
- **Configuration**: YAML-based config with environment override
- **Logging**: Zap logger with lumberjack rotation
- **Documentation**: Swagger/OpenAPI auto-generated docs at `/swagger/*`

### API Structure
All endpoints are under `/api/` with versioned routes at `/api/v1/` (with auth middleware):
- `POST /api/timelogs` - Create time log
- `GET /api/timelogs` - List time logs
- `GET /api/timelogs/:id` - Get specific time log
- `PUT /api/timelogs/:id` - Update time log
- `DELETE /api/timelogs/:id` - Delete time log

### Configuration Files
- `config.yml`: Production config
- `config-test.yml`: Test environment config
- `config-example.yml`: Template configuration

Database path, server address/port, CORS origins, and logging settings are all configurable.

### Testing Strategy
Integration tests use a separate test database and can flush data between runs. Test setup includes a fake logger implementation and proper migration handling.
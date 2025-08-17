# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a full-stack time logging application that allows users to track time entries with structured tags and remarks. It's a single-user local system with no authentication, built using:

- **Backend**: Go with Gin framework, SQLite database, and GORM ORM
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS
- **Architecture**: RESTful API with modern web frontend

## Common Development Commands

### Building
```bash
# Build backend for current platform
make build

# Build backend for Linux ARM64 (for Docker)
make build-linux

# Build frontend
make web-build

# Build and create Docker image (includes frontend)
make docker env=prod
make docker env=test
```

### Running
```bash
# Local development - Backend only
./main                # Runs on port 8083

# Local development - Full stack
./main &              # Start backend on port 8083
make web-dev          # Start frontend dev server on port 3000

# Production mode
make web-build && ./main   # Build frontend then start backend

# Docker containers
make run env=prod     # Runs on port 8083
make run env=test     # Runs on port 18083

# Manual Docker run
docker run --rm -e ENV=prod -p 8083:8083 timelog-app
docker run --rm -e ENV=test -p 18083:8083 timelog-app-test
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
All endpoints are under `/api/` with unified response format:

**TimeLog Endpoints:**
- `POST /api/timelogs` - Create time log
- `GET /api/timelogs` - List time logs (with tag info)
- `GET /api/timelogs/:id` - Get specific time log
- `PUT /api/timelogs/:id` - Update time log
- `DELETE /api/timelogs/:id` - Delete time log

**Tag Endpoints:**
- `GET /api/tags` - List all available tags
- `POST /api/tags` - Create new tag
- `GET /api/tags/:id` - Get specific tag
- `PUT /api/tags/:id` - Update tag
- `DELETE /api/tags/:id` - Delete tag

**Response Format:**
```json
{
  "data": {...},
  "message": "Operation successful",
  "status": 200
}
```

### Configuration Files
- `config.yml`: Production config (port 8083)
- `config-test.yml`: Test environment config (port 8083)
- `config-example.yml`: Template configuration

Database path, server address/port, CORS origins, and logging settings are all configurable.

### Database Schema
The application uses two main tables:

**timelogs table:**
- `id` - Primary key
- `user_id` - User identifier (default 1)
- `start_time` - When the time log started (required)
- `end_time` - When the time log ended (nullable)
- `tag_id` - Foreign key to tags table (required)
- `remark` - Additional notes/description
- Timestamps: `created_at`, `updated_at`, `deleted_at`

**tags table:**
- `id` - Primary key
- `name` - Tag display name (unique)
- `color` - Hex color code for UI display
- `description` - Tag description/purpose
- Timestamps: `created_at`, `updated_at`, `deleted_at`

**Default Tags:**
- 工作 (Work) - Red #EF4444
- 学习 (Study) - Green #10B981
- 会议 (Meeting) - Yellow #F59E0B
- 开发 (Development) - Purple #8B5CF6
- 休息 (Rest) - Gray #6B7280
- 运动 (Exercise) - Orange #F97316
- 其他 (Other) - Blue #6366F1

### Frontend Structure
Located in `web/` directory:
- **Framework**: Vue 3 with Composition API
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Icons**: Heroicons
- **HTTP Client**: Axios

### Testing Strategy
Integration tests use a separate test database and can flush data between runs. Test setup includes a fake logger implementation and proper migration handling.

### Development Ports
- Backend API: 8083
- Frontend Dev Server: 3000 (proxies API to backend)
- Production: 8083 (serves both API and static frontend files)

## Recent Major Changes

### Tag System Implementation (Latest)
- Converted tag field from free text to structured Tag table with foreign key relationship
- Added Tag CRUD API endpoints with color and description support
- Updated frontend to use tag selection dropdown instead of text input
- Implemented tag color display in UI with 7 predefined tags
- Database migration completed to preserve existing data

### API Response Format Standardization
- Unified all API responses to use consistent format: `{data, message, status}`
- Updated frontend to handle new response structure
- Improved error handling and user feedback

### Full-Stack Integration
- Vue 3 frontend with TypeScript and Tailwind CSS
- Backend serves both API (port 8083) and static frontend files
- Development proxy setup for frontend dev server (port 3000 → 8083)
- Production build integration with Go static file serving
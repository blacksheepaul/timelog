# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive full-stack time logging and task management application. Users can track time entries with structured tags, manage tasks with deadlines and estimations, and view analytical reports. It's a single-user local system with no authentication, built using:

- **Backend**: Go with Gin framework, SQLite database, and GORM ORM
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS with Vue Router
- **Architecture**: RESTful API with modern SPA frontend
- **Features**: Time logging, task management, reporting, and analytics

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

### Code Formatting
```bash
# After making modifications and completing tests, reformat the code
make fmt
```

## Architecture

### Layer Structure
- **main.go**: Application entry point with graceful shutdown
- **router/**: HTTP routing and handlers
  - `router.go`: Main router setup with middleware
  - `timelog.go`: TimeLog API endpoints (CRUD operations)
  - `task.go`: Task API endpoints (CRUD operations)
  - `middleware/`: CORS and auth middleware
- **service/**: Business logic layer
  - `timelog.go`: TimeLog business operations
  - `task.go`: Task business operations and task-timelog integration
- **model/**: Data access layer with GORM
  - `timelog.go`: TimeLog model and CRUD functions
  - `task.go`: Task model and CRUD functions
  - `tag.go`: Tag model and CRUD functions
  - `model.go`: Database connection and DAO pattern
  - `migrations/`: SQL migration files (includes task system)
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
- `POST /api/timelogs` - Create time log (with optional task association)
- `GET /api/timelogs` - List time logs (with tag and task info)
- `GET /api/timelogs/:id` - Get specific time log
- `PUT /api/timelogs/:id` - Update time log
- `DELETE /api/timelogs/:id` - Delete time log

**Task Endpoints:**
- `POST /api/tasks` - Create new task
- `GET /api/tasks` - List all tasks (with optional date filter)
- `GET /api/tasks/:id` - Get specific task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `POST /api/tasks/:id/complete` - Mark task as completed
- `POST /api/tasks/:id/incomplete` - Mark task as incomplete
- `GET /api/tasks/stats/:date` - Get task completion statistics for date

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
The application uses three main tables with foreign key relationships:

**timelogs table:**
- `id` - Primary key
- `user_id` - User identifier (default 0)
- `start_time` - When the time log started (required)
- `end_time` - When the time log ended (nullable)
- `tag_id` - Foreign key to tags table (required)
- `task_id` - Foreign key to tasks table (nullable, for task association)
- `remark` - Additional notes/description
- Timestamps: `created_at`, `updated_at`, `deleted_at`

**tasks table:**
- `id` - Primary key
- `title` - Task title (required, max 200 chars)
- `description` - Detailed task description (optional)
- `tag_id` - Foreign key to tags table (required)
- `due_date` - Task due date (required)
- `estimated_minutes` - Estimated completion time in minutes (required)
- `is_completed` - Boolean completion status (default false)
- `completed_at` - Timestamp when task was completed (nullable)
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
- 恢复 (Recovery) - Yellow #F59E0B
- 开发 (Development) - Purple #8B5CF6
- 休息 (Rest) - Gray #6B7280
- 运动 (Exercise) - Orange #F97316
- 其他 (Other) - Blue #6366F1

### Frontend Structure
Located in `web/` directory:
- **Framework**: Vue 3 with Composition API
- **Routing**: Vue Router 4 for SPA navigation
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Icons**: Heroicons
- **HTTP Client**: Axios
- **Date Utils**: date-fns

**Page Structure:**
- `src/views/Home.vue` - Dashboard with stats and quick actions
- `src/views/TimeLog.vue` - Time logging management
- `src/views/Tasks.vue` - Task management interface
- `src/views/Tags.vue` - Tag administration  
- `src/views/Statistics.vue` - Analytics and reporting

**Key Features:**
- Multi-page SPA with responsive navigation
- Real-time duration calculation with timezone handling
- Task management with deadline tracking
- Color-coded tag system
- Form validation and error handling
- Toast notification system

### Testing Strategy
Integration tests use a separate test database and can flush data between runs. Test setup includes a fake logger implementation and proper migration handling.

### Development Ports
- Backend API: 8083
- Frontend Dev Server: 3000 (proxies API to backend)
- Production: 8083 (serves both API and static frontend files)

## Recent Major Changes

### Task Management System (Latest)
- **Complete task management functionality** added to application
- New `tasks` table with title, description, tag association, due dates, and time estimates
- Task CRUD API endpoints with completion status management
- Frontend task management interface with creation, editing, and status tracking
- Task-TimeLog association support (task_id field added to timelogs table)
- Date filtering and task statistics API
- Database migrations: `000005_create_tasks_table.sql` and `000006_add_task_id_to_timelogs.sql`

### Vue Router Integration
- **Multi-page SPA architecture** implemented with Vue Router 4
- Navigation system with responsive design for desktop and mobile
- Five main pages: Dashboard, Time Logs, Tasks, Tags, Statistics
- Route-based page titles and active navigation indicators

### Enhanced Time Management Features
- **"Set to Now" functionality** for quick time entry
- Fixed timezone handling issues between frontend and backend
- Improved duration calculation with proper date-fns formatting
- Real-time time display with timezone consistency (UTC storage, local display)

### Tag System Implementation 
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

## Future Development Areas

### Planned Features
1. **Task-TimeLog Integration**: Enhanced workflow where completing tasks can optionally create time logs
2. **Automated Reporting**: Daily reports generated at 4 AM showing task completion vs. time estimates
3. **Advanced Analytics**: Better visualization of productivity patterns and time allocation
4. **Task Templates**: Reusable task templates for common work patterns
5. **Bulk Operations**: Multi-select and bulk actions for tasks and time logs
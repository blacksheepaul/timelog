# timelog

Yet another lyubishchev time management implementation.

Core functions:

- Time tracking
- LLM review based on time tracking data (via MCP)
- Task management (wip)

# How to use

# How to build and run

## Swagger Setup (Development/Testing Only)

For development and testing environments, you need to generate Swagger documentation using the `swag` tool before running `go mod tidy` or building the project:

### Install swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger documentation

```bash
swag init
```

This will generate the necessary files in the `docs` directory that are required by the router package.

**Note:** For production builds, Swagger is automatically excluded via the `prod` build tag, so you don't need to generate Swagger documentation for production deployments.

## Migrate

for example:

create new migration

```bash
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ init_xxx_table
```

forward

```bash
migrate -database "sqlite3://dev.db" --path model/migrations/ up
# or use make target (defaults to dev environment)
make migrate
# or explicitly specify environment
make migrate env=prod
make migrate env=dev
```

# Launch

```bash
# Build binary
make build env=prod
# or
make build-linux env=dev
# Build image
docker build -t timelog-app .
# Prod
docker run --rm -e ENV=prod -p 8080:8080 timelog-app
# Dev
docker run --rm -e ENV=dev -p 18080:8080 timelog-app
```

# How to Deployment

- English: [deploy.md](docs/deploy.md)
- 中文: [deploy-cn.md](docs/deploy-cn.md)

### TODO

- frontend: +manual refresh/fetch

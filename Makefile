# Makefile for timelog project

# Default values
env ?= prod
migrate_env ?= dev
BIN_NAME := main
BIN_NAME_ARM := main.arm
BIN_NAME_LINUX := main.linux

ifeq ($(env),prod)
	PORT := 8083
	DOCKER_TAG := timelog-app
	DOCKER_PORT := 8083
else ifeq ($(env),dev)
	PORT := 18083
	DOCKER_TAG := timelog-app-dev
	DOCKER_PORT := 8083
else
	$(error Unknown env: $(env))
endif

# Set DB file for migrate target
ifeq ($(env),prod)
	MIGRATE_DB_FILE := prod.db
else
	MIGRATE_DB_FILE := dev.db
endif

.PHONY: all build build-linux buildx buildx-linux docker run clean web web-build web-dev mcp-server migrate

all: build

build:
	go build -o $(BIN_NAME)

build-linux:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN_NAME_ARM)

buildx: web-build build

buildx-linux: web-build build-linux

docker: build-linux
	docker build -t $(DOCKER_TAG) .

run: docker
	docker run --rm -e ENV=$(env) -p $(PORT):$(DOCKER_PORT) $(DOCKER_TAG)

clean:
	rm -f $(BIN_NAME) $(BIN_NAME_ARM)
	rm -rf web/dist web/node_modules

# Web frontend targets
web-build:
	cd web && pnpm install && pnpm run build

web-dev:
	cd web && pnpm install && pnpm run dev

web: web-build

# MCP Server target
mcp-server:
	cd mcp && go build -o timelog-mcp-server .

# Migrate target
migrate:
	migrate -database "sqlite3://$(MIGRATE_DB_FILE)" --path model/migrations/ up

# fmt
fmt:
	go fmt ./...
	cd mcp && go fmt ./...
	cd web && npx prettier --write src/ || true

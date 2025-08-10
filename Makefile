# Makefile for timelog project

# Default values
env ?= prod
BIN_NAME := main
BIN_NAME_ARM := main.arm

ifeq ($(env),prod)
	PORT := 8080
	DOCKER_TAG := timelog-app
	DOCKER_PORT := 8080
else ifeq ($(env),test)
	PORT := 18080
	DOCKER_TAG := timelog-app-test
	DOCKER_PORT := 8080
else
	$(error Unknown env: $(env))
endif

.PHONY: all build build-linux docker run clean

all: build

build:
	go build -o $(BIN_NAME)

build-linux:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN_NAME_ARM)

docker: build-linux
	docker build -t $(DOCKER_TAG) .

run: docker
	docker run --rm -e ENV=$(env) -p $(PORT):$(DOCKER_PORT) $(DOCKER_TAG)

clean:
	rm -f $(BIN_NAME) $(BIN_NAME_ARM)

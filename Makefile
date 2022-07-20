# Base Go commands.
GO_CMD     := go
GO_FMT     := $(GO_CMD) fmt
GO_INSTALL := $(GO_CMD) install
GO_CLEAN   := $(GO_CMD) clean
GO_BUILD   := $(GO_CMD) build -mod vendor

# Coverage output.
COVER_OUT := cover.out

# Base swagger commands.
SWAG     := swag
SWAG_GEN := $(SWAG) init

# Base golangci-lint commands.
GCL_CMD := golangci-lint
GCL_RUN := $(GCL_CMD) run

# Base protoc commands.
PROTOC := protoc

# Project executable file, and its binary.
CMD_PATH    := ./cmd/akatsuki
BINARY_NAME := akatsuki

# Default makefile target.
.DEFAULT_GOAL := run

# Standarize go coding style for the whole project.
.PHONY: fmt
fmt:
	@$(GO_FMT) ./...

# Lint go source code.
.PHONY: lint
lint: fmt
	@$(GCL_RUN) -D errcheck --timeout 5m

# Clean project binary, test, and coverage file.
.PHONY: clean
clean:
	@$(GO_CLEAN) ./...

# Generate swagger docs.
.PHONY: swagger
swagger:
	@$(SWAG_GEN) --parseVendor -g cmd/akatsuki/main.go -o ./docs

# Generate grpc proto buff.
.PHONY: grpc
grpc:
	@$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/delivery/grpc/schema/*.proto

# Install library.
.PHONY: install
install:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2
	@$(GCL_CMD) version
	@$(GO_INSTALL) github.com/swaggo/swag/cmd/swag@v1.8.3
	@$(SWAG) -v

# Build the project executable binary.
.PHONY: build
build: clean fmt
	@cd $(CMD_PATH); \
	$(GO_BUILD) -o $(BINARY_NAME) -v .

# Build and run the binary.
.PHONY: run
run: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) server

# Build and migrate database.
.PHONY: migrate
migrate: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) migrate

# Build and run message consumer.
.PHONY: consumer
consumer: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) consumer

# Build and run cron update anime.
.PHONY: cron-update
cron-update: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) cron update

# Build and run cron fill missing anime.
.PHONY: cron-fill
cron-fill: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) cron fill

# Docker base command.
DOCKER_CMD   := docker
DOCKER_IMAGE := $(DOCKER_CMD) image

# Docker-compose base command and docker-compose.yml path.
COMPOSE_CMD                 := docker-compose
COMPOSE_BUILD               := deployment/build.yml
COMPOSE_API                 := deployment/api.yml
COMPOSE_CONSUMER            := deployment/consumer.yml
COMPOSE_CRON_UPDATE         := deployment/cron-update.yml
COMPOSE_CRON_FILL           := deployment/cron-fill.yml
COMPOSE_CRON_CALLBACK_RETRY := deployment/cron-callback-retry.yml
COMPOSE_MIGRATE             := deployment/migrate.yml
COMPOSE_LINT                := deployment/lint.yml

# Build docker images and container for the project
# then delete builder image.
.PHONY: docker-build
docker-build:
	@$(COMPOSE_CMD) -f $(COMPOSE_BUILD) build
	@$(DOCKER_IMAGE) prune -f --filter label=stage=akatsuki_builder

# Start built docker containers for api.
.PHONY: docker-api
docker-api:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p akatsuki-api up -d
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p akatsuki-api logs --follow --tail 20

# Start built docker containers for consumer.
.PHONY: docker-consumer
docker-consumer:
	@$(COMPOSE_CMD) -f $(COMPOSE_CONSUMER) -p akatsuki-consumer up -d
	@$(COMPOSE_CMD) -f $(COMPOSE_CONSUMER) -p akatsuki-consumer logs --follow --tail 20

# Start built docker containers for cron update anime.
.PHONY: docker-cron-update
docker-cron-update:
	@$(COMPOSE_CMD) -f $(COMPOSE_CRON_UPDATE) -p akatsuki-cron-update up

# Start built docker containers for cron fill missing anime.
.PHONY: docker-cron-fill
docker-cron-fill:
	@$(COMPOSE_CMD) -f $(COMPOSE_CRON_FILL) -p akatsuki-cron-fill up

# Start built docker containers for migrate.
.PHONY: docker-migrate
docker-migrate:
	@$(COMPOSE_CMD) -f $(COMPOSE_MIGRATE) -p akatsuki-migrate up

# Start docker to run lint check.
.PHONY: docker-lint
docker-lint:
	@$(COMPOSE_CMD) -f $(COMPOSE_LINT) -p akatsuki-lint run --rm akatsuki-lint $(GCL_RUN) -D errcheck --timeout 5m

# Update docker containers.
.PHONY: docker-update
docker-update:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p akatsuki-api up -d
	@$(COMPOSE_CMD) -f $(COMPOSE_CONSUMER) -p akatsuki-consumer up -d
	@$(DOCKER_IMAGE) prune -f --filter label=stage=akatsuki_binary

# Stop running docker containers.
.PHONY: docker-stop
docker-stop:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p akatsuki-api stop
	@$(COMPOSE_CMD) -f $(COMPOSE_CONSUMER) -p akatsuki-consumer stop
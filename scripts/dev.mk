# Go command.
GO ?= go

# Tool versions.
GCL_VERSION  := v2.15.0
SWAG_VERSION := v1.16.4

# Test coverage output.
COVER_OUT := cover.out

# Project paths.
CMD_PATH    := ./cmd/akatsuki
BINARY_NAME := akatsuki
BINARY_PATH := $(CMD_PATH)/$(BINARY_NAME)

.PHONY: fmt clean lint swagger grpc test install-linter install-swagger install-grpc \
	build run migrate consumer cron-update cron-fill

##@ Tools

fmt: ## Standarize go coding style for the whole project.
	@$(GO) fmt ./...

clean: ## Clean project binary, test, and coverage file.
	@$(GO) clean ./...
	@rm -f $(COVER_OUT)

lint: fmt ## Lint go source code.
	@golangci-lint run -D errcheck --timeout 5m

swagger: ## Generate swagger docs.
	@swag init -g $(CMD_PATH)/main.go -o ./docs

grpc: ## Generate grpc proto buf.
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative internal/delivery/grpc/schema/*.proto

test: ## Run test.
	@$(GO) test -v -cover -race -timeout=1m -coverprofile=$(COVER_OUT) ./...
	@$(GO) tool cover -func=$(COVER_OUT)
	@rm -f $(COVER_OUT)

install-linter: ## Install golangci-lint.
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GCL_VERSION)

install-swagger: ## Install swag.
	@$(GO) install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

install-grpc: ## Install protoc.
	@echo "please go to https://grpc.io/docs/languages/go/quickstart/"

##@ Build & Run

build: clean fmt ## Build the project executable binary. 
	@$(GO) build -o $(BINARY_PATH) -v $(CMD_PATH)

migrate: build ## Build and migrate database.
	@$(BINARY_PATH) migrate

run: build ## Build and run the server.
	@$(BINARY_PATH) server

consumer: build ## Build and run consumer.
	@$(BINARY_PATH) consumer

cron-update: build ## Build and run cron update.
	@$(BINARY_PATH) cron update

cron-fill: build ## Build and run cron fill.
	@$(BINARY_PATH) cron fill
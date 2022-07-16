# Base Go commands.
GO_CMD     := go
GO_FMT     := $(GO_CMD) fmt
GO_GET     := $(GO_CMD) get
GO_INSTALL := $(GO_CMD) install
GO_MOD     := $(GO_CMD) mod
GO_CLEAN   := $(GO_CMD) clean
GO_BUILD   := $(GO_CMD) build -mod vendor
GO_RUN     := $(GO_CMD) run -mod vendor
GO_TEST    := $(GO_CMD) test
GO_COVER   := $(GO_CMD) tool cover

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

# Run test.
.PHONY: test
test:
	@$(GO_TEST) ./... -coverprofile $(COVER_OUT)
	@$(GO_COVER) -func $(COVER_OUT)

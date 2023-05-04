envfile ?= .env.dist
-include $(envfile)
ifneq ("$(wildcard $(envfile))","")
	export $(shell sed 's/=.*//' $(envfile))
endif

GOLANGCI_VERSION:=1.52.2
PROJECT_NAME:=cosmos-grpc-forwarder
SERVER_NAME:=grpc-server
CLIENT_NAME:=grpc-client
GOPATH_BIN:=$(shell go env GOPATH)/bin
DOCKER := $(shell which docker)
PROTO_DOCKER_VERSION=0.12.1
PROTO_DOCKER_IMAGE_NAME=ghcr.io/cosmos/proto-builder:$(PROTO_DOCKER_VERSION)
PROTO_IMAGE_EXEC=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_DOCKER_IMAGE_NAME)

.PHONY: install
install:
	# Install golangci-lint for go code linting.
	curl -sSfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | \
		sh -s -- -b ${GOPATH_BIN} v${GOLANGCI_VERSION}


.PHONY: all
all: clean init lint test build-server build-client

.PHONY: init
init:
	@cp .env.dist .env
	@cp .env.test.dist .env.test

.PHONY: lint
lint:
	@echo ">>> Performing golang code linting.."
	golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@echo ">>> Running Unit Tests..."
	go test -v -race ./...

.PHONY: cover-test
cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -v -race ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: build-server
build-server:
	@echo ">>> Building ${PROJECT_NAME} gRPC server..."
	go build -o bin/${SERVER_NAME} cmd/${SERVER_NAME}/main.go

.PHONY: build-client
build-client:
	@echo ">>> Building ${PROJECT_NAME} gRPC client..."
	go build -o bin/${CLIENT_NAME} cmd/${CLIENT_NAME}/main.go

.PHONY: run-server
run-server:
	@echo ">>> Running ${PROJECT_NAME} gRPC server..."
	@go run ./cmd/${SERVER_NAME}/main.go

.PHONY: run-client
run-client:
	@echo ">>> Running ${PROJECT_NAME} gRPC client..."
	@go run ./cmd/${CLIENT_NAME}/main.go

.PHONY: clean
clean:
	@echo ">>> Removing old binaries and env files..."
	@rm -rf bin/*
	@rm -rf .env

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(PROTO_IMAGE_EXEC) buf generate

proto-format:
	@$(PROTO_IMAGE_EXEC) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(PROTO_IMAGE_EXEC) buf lint --error-format=json

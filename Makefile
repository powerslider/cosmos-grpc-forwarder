envfile ?= .env.dist
-include $(envfile)
ifneq ("$(wildcard $(envfile))","")
	export $(shell sed 's/=.*//' $(envfile))
endif

BUF_VERSION:=1.17.0
GOLANGCI_VERSION:=1.52.2
PROJECT_NAME:=cosmos-grpc-forwarder
GOPATH_BIN:=$(shell go env GOPATH)/bin

.PHONY: install
install:
	# Install golangci-lint for go code linting.
	curl -sSfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | \
		sh -s -- -b ${GOPATH_BIN} v${GOLANGCI_VERSION}

	# Install buf tool for protobuf stub generation, linting, etc.
	go install github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}

	git clone https://github.com/cosmos/gogoproto.git; \
      cd gogoproto; \
      go mod download; \
      make install


.PHONY: all
all: clean init lint test build-server

.PHONY: init
init:
	@cp .env.dist .env

.PHONY: lint
lint:
	@echo ">>> Performing golang code linting.."
	golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@echo ">>> Running Unit Tests..."
	go test -race ./...

.PHONY: cover-test
cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: build-server
build-server:
	@echo ">>> Building ${PROJECT_NAME} API server..."
	go build -o bin/server cmd/${PROJECT_NAME}/main.go

.PHONY: run-server
run-server:
	@echo ">>> Running ${PROJECT_NAME} API server..."
	@go run ./cmd/${PROJECT_NAME}/main.go

.PHONY: docs
docs:
	@echo ">>> Generate Swagger API Documentation..."
	swag init --generalInfo cmd/${PROJECT_NAME}/main.go

.PHONY: clean
clean:
	@echo ">>> Removing old binaries and env files..."
	@rm -rf bin/*
	@rm -rf .env


DOCKER := $(shell which docker)
protoVer=0.12.1
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) buf generate

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	$(MAKE) update-swagger-docs

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main
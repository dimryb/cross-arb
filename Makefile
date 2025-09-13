#Makefile
NAME := "cross-arb"

GIT_HASH := $(shell git log --format="%h" -n 1)

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	BIN := ./bin/$(NAME)
	BIN_CLIENT := ./bin/gRPC-client
    DATE_CMD = date -u +'%Y-%m-%dT%H:%M:%S'
    GO_PATH := $(shell go env GOPATH)
else #windows
	BIN := ./bin/$(NAME).exe
	BIN_CLIENT := ./bin/gRPC-client.exe
    DATE_CMD = powershell.exe -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'"
    GO_PATH := $(shell go env GOPATH | tr '\\' '/')
endif

LDFLAGS := -X main.release="develop" \
    -X main.buildDate=$(shell $(DATE_CMD)) \
    -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/app

.PHONY: run
run: build
	$(BIN) -config ./configs/config.yaml

.PHONY: version
version: build
	$(BIN) version

.PHONY: build-client
build-client:
	go build -v -o $(BIN_CLIENT) ./cmd/client

.PHONY: run-client
run-client: build-client
	$(BIN_CLIENT)

.PHONY: test
test:
	go test -race ./internal/...

.PHONY: install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	sh -s -- -b $(GO_PATH)/bin v1.64.8

.PHONY: lint
lint: install-lint-deps
	golangci-lint run --config golangci.yml ./...

.PHONY: install-arch-deps
install-arch-deps:
	(which go-arch-lint > /dev/null) || \
	GO111MODULE=on go install github.com/fe3dback/go-arch-lint@v1.12.0

.PHONY: arch
arch: install-lint-deps
	go-arch-lint check

.PHONY: graph
graph: install-lint-deps
	go-arch-lint graph

.PHONY: generate
generate:
	protoc \
		-I proto \
		--go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		proto/*.proto

.PHONY: generate-mocks
generate-mocks:
	go generate ./...

# Установи grpcurl
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
.PHONY: grpc-subscribe
grpc-subscribe:
	grpcurl -plaintext localhost:9090 ticker.TickerService/Subscribe

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

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/app

run: build
	$(BIN) -config ./configs/config.yaml

version: build
	$(BIN) version

build-client:
	go build -v -o $(BIN_CLIENT) ./cmd/client

run-client: build-client
	$(BIN_CLIENT)

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	sh -s -- -b $(GO_PATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run --config golangci.yml ./...

arch:
	go-arch-lint check

graph:
	go-arch-lint graph

generate:
	protoc \
		-I proto \
		--go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		proto/*.proto

generate-mocks:
	go generate ./...

# Установи grpcurl
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
grpc-subscribe:
	grpcurl -plaintext localhost:9090 ticker.TickerService/Subscribe

.PHONY: build run version test install-lint-deps lint arch graph generate generate-mocks grpc-subscribe

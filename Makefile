BINARY_NAME := go-serve
VERSION := $(shell git describe --tags)
TIME := $(shell date +%Y-%m-%dT%T%z)
BUILD := $(shell git rev-parse --short HEAD)
DIST_FOLDER := ./dist
BINARY_OUTPUT := $(DIST_FOLDER)/$(BINARY_NAME)
LDFLAGS=-ldflags "-s -w \
		-X=github.com/eloylp/go-serve/server.Name=$(BINARY_NAME) \
		-X=github.com/eloylp/go-serve/server.Version=$(VERSION) \
		-X=github.com/eloylp/go-serve/server.Build=$(BUILD) \
		-X=github.com/eloylp/go-serve/server.BuildTime=$(TIME)"
FLAGS=-trimpath

.DEFAULT_GOAL := build

lint:
	golangci-lint run -v
lint-fix:
	golangci-lint run -v --fix

all: lint test build

test: test-unit test-integration test-racy

test-unit:
	go test -race -v --tags="unit" ./...
test-integration:
	go test -race -v --tags="integration" ./...
test-racy:
	go test -race -v --tags="racy" ./...
build:
	mkdir -p $(DIST_FOLDER)
	CGO_ENABLED=0 go build $(FLAGS) $(LDFLAGS) -o $(BINARY_OUTPUT) ./cmd/server
	@echo "Binary output at $(BINARY_OUTPUT)"
clean:
	rm -rf $(DIST_FOLDER)
BINARY_NAME := go-serve
VERSION := $(shell git describe --tags)
TIME := $(shell date +%Y-%m-%dT%T%z)
BUILD := $(shell git rev-parse --short HEAD)
DIST_FOLDER := ./dist
BINARY_OUTPUT := $(DIST_FOLDER)/$(BINARY_NAME)
LDFLAGS=-ldflags "-s -w \
		-X=go.eloylp.dev/go-serve/server.Name=$(BINARY_NAME) \
		-X=go.eloylp.dev/go-serve/server.Version=$(VERSION) \
		-X=go.eloylp.dev/go-serve/server.Build=$(BUILD) \
		-X=go.eloylp.dev/go-serve/server.BuildTime=$(TIME)"
FLAGS=-trimpath
TAGS=-tags timetzdata

.DEFAULT_GOAL := build

lint:
	golangci-lint run --build-tags unit,integration,racy -v
lint-fix:
	golangci-lint run --build-tags unit,integration,racy -v --fix

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
	CGO_ENABLED=0 go build $(FLAGS) $(LDFLAGS) $(TAGS) -o $(BINARY_OUTPUT) ./cmd/server
	@echo "Binary output at $(BINARY_OUTPUT)"
clean:
	rm -rf $(DIST_FOLDER)
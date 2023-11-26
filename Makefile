#!/usr/bin/env make -f

PROJECTNAME := "ssh-keys"
BUILD := $(shell git rev-parse --short HEAD)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -w -X=main.build=$(BUILD)"

all: test lint build

## test: Run the go test command.
.PHONY: test
test:
	go test -v ./...

## lint: Run the linting.
.PHONY: lint
lint:
	golangci-lint run ./...
	staticcheck ./...

## build: Compile the binary.
.PHONY: build
build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(PROJECTNAME) cmd/example-gorilla-rest-api/main.go

## clean: Cleanup binary.
clean:
	@rm -f bin/$(PROJECTNAME)

## help: Show this message.
.PHONY: help
help: Makefile
	@echo "Available targets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

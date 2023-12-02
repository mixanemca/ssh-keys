#!/usr/bin/env make -f

PROJECTNAME := "ssh-keys"
BUILD := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --abbrev=0 --tags)
DATE := $(shell date)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-s -w -X 'github.com/version-go/ldflags.buildVersion=$(VERSION)' -X 'github.com/version-go/ldflags.buildHash=$(BUILD)' -X 'github.com/version-go/ldflags.buildTime=$(DATE)'"

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
	go build $(LDFLAGS) -o bin/$(PROJECTNAME) main.go

## install: Install the binary.
.PHONY: install
install:
	go install $(LDFLAGS) .

## clean: Cleanup binary.
clean:
	@rm -f bin/$(PROJECTNAME)

## help: Show this message.
.PHONY: help
help: Makefile
	@echo "Available targets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

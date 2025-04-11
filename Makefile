.PHONY: all build dev run clean test deps

BINARY_NAME=main

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/tmp

all: build

build:
	go build -o $(GOBIN)/$(BINARY_NAME) $(GOBASE)/examples/main/main.go

dev:
	air

run:
	go run $(GOBASE)/examples/main/main.go

clean:
	rm -rf $(GOBIN)

test:
	go test -v ./...

deps:
	go mod tidy

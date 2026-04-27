.PHONY: build install test fmt vet tidy

BIN := bin/aikido
PKG := github.com/xandervr/aikido-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X $(PKG)/internal/version.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/aikido

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/aikido

test:
	go test ./... -count=1 -race -cover

fmt:
	gofmt -s -w .

vet:
	go vet ./...

tidy:
	go mod tidy

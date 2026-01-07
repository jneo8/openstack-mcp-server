SHELL := /bin/bash
GO_VERSION ?= "1.24.2"
BINARY_NAME := mcp-server
BUILD_DIR := ./bin
CMD_DIR := ./cmd/mcp-server

##@ Development

run:  ## Run the application
	go run $(CMD_DIR) serve

fmt:  ## Format Go code
	go fmt ./...

vet:  ## Run go vet
	go vet ./...

tidy:  ## Tidy go modules
	go mod tidy

.PHONY: run fmt vet tidy

##@ Testing

test:  ## Run tests
	go test ./... -v

test-coverage:  ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-race:  ## Run tests with race detector
	go test ./... -race

.PHONY: test test-coverage test-race

##@ Linting
# Install golangci-lint: https://golangci-lint.run/usage/install/
# Recommended: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

lint:  ## Run golangci-lint
	golangci-lint run ./...

lint-fix:  ## Run golangci-lint with auto-fix
	golangci-lint run ./... --fix

lint-verbose:  ## Run golangci-lint with verbose output
	golangci-lint run ./... -v

staticcheck:  ## Run staticcheck
	staticcheck ./...

.PHONY: lint lint-fix lint-verbose staticcheck

##@ Quality

check: fmt vet lint test  ## Run all checks (fmt, vet, lint, test)

pre-commit: check  ## Run pre-commit checks

.PHONY: check pre-commit

##@ Build

VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

build:  ## Build the binary
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="$(LDFLAGS)" \
		-a -installsuffix cgo \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(CMD_DIR)

build-local:  ## Build the binary for local OS
	@mkdir -p $(BUILD_DIR)
	go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(CMD_DIR)

clean:  ## Clean build artifacts
	rm -rf $(BUILD_DIR)

.PHONY: build build-local clean

##@ Docker

DOCKER_REGISTRY ?= ghcr.io
DOCKER_REPO ?= jneo8/openstack-mcp-server
DOCKER_TAG ?= $(VERSION)
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_REPO):$(DOCKER_TAG)

docker-build:  ## Build Docker image
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_IMAGE) \
		.

docker-build-latest:  ## Build Docker image with latest tag
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_REGISTRY)/$(DOCKER_REPO):latest \
		-t $(DOCKER_IMAGE) \
		.

.PHONY: docker-build docker-build-latest

##@ Help

.PHONY: help

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

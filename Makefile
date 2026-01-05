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

##@ Help

.PHONY: help

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

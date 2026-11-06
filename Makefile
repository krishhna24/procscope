BINARY := procscope
PKG := github.com/krishhna24/procscope
BIN_DIR := bin

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X $(PKG)/internal/version.Version=$(VERSION)

FUZZTIME ?= 30s

# Pinned linter version. CI (guide 05) uses the same one.
GOLANGCI_LINT_VERSION := v2.15.0
GOVULNCHECK_VERSION := v1.9.0


.DEFAULT_GOAL := help
.PHONY: help build test test-race lint fmt tools fuzz-smoke bench integration clean

help: ## List the available targets
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf " %-12s %s\n", $$1, $$2}'

build: ## Build bin/procscope with version info
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) ./cmd/procscope

test: ## Run unit tests
	go test ./...

test-race: ## Run unit tests with the race detector
	go test -race ./...

lint: ## Run golangci-lint (gofmt, vet, staticcheck, errcheck, errorlint, misspell)
	@command -v golangci-lint >/dev/null || { echo "golangci-lint not found: run 'make tools'"; exit 1; }
	golangci-lint run ./...

fmt: ## Format every Go file in place
	golangci-lint fmt ./...

tools: ## Install the pinned golangci-lint into $(go env GOPATH)/bin
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

fuzz-smoke: ## Run every fuzz target for FUZZTIME (targets arrive in Phase 2)
	@echo "no fuzz targets yet"

bench: ## Run benchmarks with allocation counts
	go test -run '^$$' -bench . -benchmem ./...

integration: ## Run tests that need a live /proc
	go test -tags integration ./...

clean: ## Remove build output
	rm -rf $(BIN_DIR)

vuln: ## Report known vulnerabilities reachable from the code
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

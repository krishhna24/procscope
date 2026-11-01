BINARY := procscope
PKG := github.com/krishhna24/procscope
BIN_DIR := bin

VERSION ?= $(shell git describe --tag --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X $(PKG)/internal/version.Version=$(VERSION)

FUZZTIME ?= 30s

.DEFAULT_GOAL := help
.PHONY: help build test test-race lint fuzz-smoke bench integration clean

help: ## List the available targets
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf " %-12s %s\n", $$1, $$2}'

build: ## Build bin/procscope with version info
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) ./cmd/procscope

test: ## Run unit tests
	go test ./...

test-race: ## Run unit tests with the race detector
		go test -race ./...

lint: ## Check formatting and run go vet
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; echo "$$unformatted"; exit 1; \
	fi
	go vet ./...

fuzz-smoke: ## Run every fuzz target for FUZZTIME (targets arrive in Phase 2)
	@echo "no fuzz targets yet"

bench: ## Run benchmarks with allocation counts
	go test -run '^$$' -bench . -benchmem ./...

integration: ## Run tests that need a live /proc
	go test -tags integration ./...

clean: ## Remove build output
	rm -rf $(BIN_DIR)

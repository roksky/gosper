GO ?= go
VERSION ?= dev
BIN_DIR ?= bin
BIN := $(BIN_DIR)/gosper
PKG := ./...
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*' -not -path './bin/*')
LDFLAGS := -X github.com/roksky/gosper/report.version=$(VERSION)

.PHONY: all help fmt fmt-check vet test test-race coverage build lint check clean

all: check

help:
	@echo "Available targets:"
	@echo "  fmt         Format Go sources with gofmt"
	@echo "  fmt-check   Verify Go sources are gofmt-clean"
	@echo "  vet         Run go vet"
	@echo "  test        Run unit tests"
	@echo "  test-race   Run unit tests with the race detector"
	@echo "  coverage    Generate a coverage profile and summary"
	@echo "  build       Build the CLI into $(BIN)"
	@echo "  lint        Run golangci-lint if installed"
	@echo "  check       Run fmt-check, vet and tests"
	@echo "  clean       Remove local build and coverage outputs"

fmt:
	gofmt -w $(GO_FILES)

fmt-check:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt found unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	$(GO) vet $(PKG)

test:
	$(GO) test $(PKG)

test-race:
	$(GO) test -race $(PKG)

coverage:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/gosper

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed; skipping."; \
		echo "Install it with:"; \
		echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		echo "See docs/development.md for details."; \
	fi

check: fmt-check vet test

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html coverage.txt

# agentic-obs Makefile
# Build automation for the OBS MCP server

# Build metadata
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD   = go
GOBUILD = $(GOCMD) build
GOTEST  = $(GOCMD) test
GOVET   = $(GOCMD) vet
GOFMT   = gofmt
GOMOD   = $(GOCMD) mod

# Binary name
BINARY_NAME = agentic-obs
BINARY_WINDOWS = $(BINARY_NAME).exe

# Ldflags for version injection
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
.PHONY: all
all: test build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for Windows specifically
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) .

# Build for all platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) .

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run tests with race detector
.PHONY: test-race
test-race:
	$(GOTEST) -v -race ./...

# Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

# Check formatting
.PHONY: fmt-check
fmt-check:
	@test -z "$$($(GOFMT) -l .)" || (echo "Files need formatting:" && $(GOFMT) -l . && exit 1)

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) -w .

# Run all linting
.PHONY: lint
lint: vet fmt-check
	@echo "Lint passed"

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Verify dependencies
.PHONY: verify
verify:
	$(GOMOD) verify

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(BINARY_WINDOWS)
	rm -f $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-darwin-arm64
	rm -f coverage.out coverage.html
	rm -rf dist/

# Run the application in MCP mode
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run the application in TUI mode
.PHONY: run-tui
run-tui: build
	./$(BINARY_NAME) --tui

# Install locally
.PHONY: install
install:
	$(GOCMD) install $(LDFLAGS) .

# Release with goreleaser (requires goreleaser)
.PHONY: release
release:
	goreleaser release --clean

# Snapshot release (for testing, no publish)
.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean

# Show version info
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

# CI targets
.PHONY: ci
ci: deps lint test build

# Quick check (fast iteration)
.PHONY: check
check: fmt vet test

# Help
.PHONY: help
help:
	@echo "agentic-obs Makefile targets:"
	@echo ""
	@echo "Build:"
	@echo "  build          - Build binary for current platform"
	@echo "  build-windows  - Build Windows binary"
	@echo "  build-all      - Build for all platforms (linux, darwin, windows)"
	@echo "  install        - Install to GOPATH/bin"
	@echo ""
	@echo "Test:"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detector"
	@echo ""
	@echo "Lint:"
	@echo "  lint           - Run all linters (vet, fmt-check)"
	@echo "  vet            - Run go vet"
	@echo "  fmt            - Format code"
	@echo "  fmt-check      - Check formatting"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy go.mod"
	@echo "  verify         - Verify dependencies"
	@echo ""
	@echo "Release:"
	@echo "  release          - Create release with goreleaser"
	@echo "  release-snapshot - Test release without publishing"
	@echo ""
	@echo "Run:"
	@echo "  run            - Build and run in MCP mode"
	@echo "  run-tui        - Build and run in TUI mode"
	@echo ""
	@echo "Utility:"
	@echo "  clean          - Remove build artifacts"
	@echo "  version        - Show version info"
	@echo "  ci             - Run CI pipeline (deps, lint, test, build)"
	@echo "  check          - Quick check (fmt, vet, test)"
	@echo "  help           - Show this help"

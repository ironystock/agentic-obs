# Building agentic-obs

This document covers building, testing, and releasing agentic-obs.

## Prerequisites

- **Go 1.21+** (tested with Go 1.25)
- **Git** (for version injection)
- **Make** (optional, for automation)
- **GoReleaser** (optional, for releases)

## Quick Start

```bash
# Clone the repository
git clone https://github.com/ironystock/agentic-obs.git
cd agentic-obs

# Build
go build -o agentic-obs .

# Run
./agentic-obs --help
```

## Build Commands

### Using Make (Recommended)

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Full CI pipeline
make ci

# See all targets
make help
```

### Using Go Directly

```bash
# Simple build
go build -o agentic-obs .

# Build with version injection
go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o agentic-obs .

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o agentic-obs.exe .

# Cross-compile for macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o agentic-obs-darwin-arm64 .
```

## Version Injection

The build system injects version information at compile time using Go's `-ldflags`:

| Variable | Description | Example |
|----------|-------------|---------|
| `main.version` | Semantic version | `1.0.0`, `1.0.0-beta.1` |
| `main.commit` | Git commit hash | `abc1234` |
| `main.date` | Build timestamp (ISO 8601) | `2025-01-15T10:30:00Z` |

### Check Version

```bash
./agentic-obs --version
# Output: agentic-obs 1.0.0

# Development build shows:
# agentic-obs dev (commit: none, built: unknown)
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...

# Run specific package tests
go test -v ./internal/mcp/...
go test -v ./internal/storage/...
```

## Cross-Platform Builds

### Supported Platforms

| OS | Architecture | Binary Name |
|----|--------------|-------------|
| Linux | amd64 | `agentic-obs-linux-amd64` |
| Linux | arm64 | `agentic-obs-linux-arm64` |
| macOS | amd64 (Intel) | `agentic-obs-darwin-amd64` |
| macOS | arm64 (Apple Silicon) | `agentic-obs-darwin-arm64` |
| Windows | amd64 | `agentic-obs.exe` |
| Windows | arm64 | `agentic-obs-windows-arm64.exe` |

### Build All Platforms

```bash
make build-all
```

Or manually:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o agentic-obs-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o agentic-obs-linux-arm64 .

# macOS
GOOS=darwin GOARCH=amd64 go build -o agentic-obs-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o agentic-obs-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o agentic-obs.exe .
GOOS=windows GOARCH=arm64 go build -o agentic-obs-windows-arm64.exe .
```

## Releases

### Using GoReleaser

GoReleaser automates the release process including:
- Cross-platform builds
- Archive creation (tar.gz for Unix, zip for Windows)
- Changelog generation
- GitHub release publishing

```bash
# Test release (no publish)
make release-snapshot

# Create release (requires GITHUB_TOKEN)
make release
```

### Manual Release

1. Tag the version:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. Build all platforms:
   ```bash
   make build-all
   ```

3. Create archives and upload to GitHub Releases.

## Configuration Files

| File | Purpose |
|------|---------|
| `Makefile` | Build automation |
| `.goreleaser.yml` | Release configuration |
| `version.go` | Version variables |
| `go.mod` | Go module definition |

## Environment Variables

Build-time environment variables:

| Variable | Description |
|----------|-------------|
| `VERSION` | Override version string |
| `COMMIT` | Override commit hash |
| `DATE` | Override build date |

Example:
```bash
VERSION=1.0.0-custom make build
```

## Continuous Integration

The project uses GitHub Actions for CI. See `.github/workflows/go.yml`.

### CI Pipeline

1. **Lint**: `go vet`, `gofmt`
2. **Test**: `go test ./...`
3. **Build**: Verify compilation

### Future: Automated Releases

GitHub Actions release workflow is planned. See [ROADMAP.md](../design/ROADMAP.md) for details.

## Troubleshooting

### CGO Issues

agentic-obs uses pure Go SQLite (`modernc.org/sqlite`), so CGO is not required:

```bash
CGO_ENABLED=0 go build -o agentic-obs .
```

### Version Shows "dev"

If `--version` shows `dev`, the build wasn't done with ldflags. Use:

```bash
make build  # Uses ldflags automatically
```

### Cross-Compilation Fails

Ensure you have Go 1.21+ which supports all target platforms natively.

## Related Documentation

- [README.md](../README.md) - User documentation
- [CLAUDE.md](../CLAUDE.md) - AI assistant context
- [ARCHITECTURE.md](../design/ARCHITECTURE.md) - System design
- [ROADMAP.md](../design/ROADMAP.md) - Future plans

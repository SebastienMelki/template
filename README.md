# {{PROJECT}}

{{PROJECT_DESCRIPTION}}

## Architecture

This project follows a **contract-first, domain-driven** architecture:

```
├── contracts/    # Protobuf API definitions (source of truth)
├── backend/      # Go service with DDD/hexagonal architecture
├── frontend/     # React + TypeScript web application
└── infra/        # Pulumi infrastructure-as-code (Go)
```

APIs are defined in Protocol Buffers and code is generated using [sebuf](https://github.com/SebastienMelki/sebuf).

## Prerequisites

| Tool | Install |
|------|---------|
| [Go 1.24+](https://go.dev) | `brew install go` |
| [Buf CLI](https://buf.build) | `brew install bufbuild/buf/buf` |
| [sebuf](https://github.com/SebastienMelki/sebuf) | `go install github.com/SebastienMelki/sebuf/cmd/...@latest` |
| [sqlc](https://sqlc.dev) | `brew install sqlc` |
| [golang-migrate](https://github.com/golang-migrate/migrate) | `brew install golang-migrate` |
| [golangci-lint](https://golangci-lint.run) | `brew install golangci-lint` |
| [Node.js 20+](https://nodejs.org) | `brew install node` |
| [Pulumi](https://pulumi.com) | `brew install pulumi` |

## Quick Start

```bash
# Generate code from proto definitions
make generate

# Run backend (in-memory database, no external deps)
make backend-run

# Run frontend dev server
make frontend-dev

# Run everything
make dev
```

## Development Commands

```bash
make generate       # Generate all code (proto + sqlc)
make build          # Build everything
make test           # Run all tests
make lint           # Lint everything
make clean          # Clean generated files and build artifacts
make backend-run    # Run backend with in-memory DB
make backend-test   # Run backend tests
make frontend-dev   # Start frontend dev server
make frontend-build # Build frontend for production
make infra-preview  # Preview infrastructure changes
make infra-up       # Deploy infrastructure
```

## Project Structure

- **contracts/** — Protocol Buffer API definitions. Source of truth for all APIs.
- **backend/** — Go backend with DDD/hexagonal architecture. Each business domain is a self-contained module.
- **frontend/** — React + Vite + TypeScript. Uses generated TypeScript clients from contracts.
- **infra/** — Pulumi infrastructure written in Go. Manages cloud resources.
- **api-docs/** — Generated OpenAPI specifications.

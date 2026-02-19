# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A **GitHub template repository** for bootstrapping full-stack projects. When a developer creates a new repo from this template and opens Claude Code, guide them through setup.

The template uses a **contract-first, domain-driven** architecture:

```
├── contracts/    # Protobuf API definitions (source of truth)
├── backend/      # Go service with DDD/hexagonal architecture
├── frontend/     # Web frontend (scaffolded during setup)
└── infra/        # Pulumi infrastructure-as-code (empty scaffold)
```

---

## FIRST-TIME SETUP

When a developer first opens this repo, walk them through setup. Keep it fast — ask only what's needed, then scaffold.

### Step 1: Project Identity

Ask:
1. **Project name** — for Go module path, package names, Pulumi stack
2. **GitHub owner** — e.g. `github.com/acme`
3. **Brief description** — one line

### Step 2: Domain Modules

Ask:
1. **Business domains** — e.g. `user`, `order`, `product`
   - For each: 2-3 initial RPC methods
2. **Shared identifiers** — e.g. `UserID`, `OrderID`

### Step 3: Scaffold

**IMPORTANT:** The placeholders below use double-brace syntax in the actual template files. When replacing, search for the literal strings including the braces.

1. Replace all placeholders: `ORG` (GitHub owner), `PROJECT` (project name), `PROJECT_DESCRIPTION` (one-liner)
2. Create proto files for each domain (service + messages)
3. Create backend module skeleton for each domain (copy the `example` module pattern)
4. Run `make generate` (sqlc generate, then buf dep update + buf generate)
5. Run `cd backend && go mod tidy` (after generation so all imports resolve)
6. Verify: `buf lint`, `cd backend && go build ./...`, `cd backend && go test ./...`
7. Commit

---

## Architecture

### Contract-First Development

APIs are defined in Protocol Buffers under `contracts/`. Code is generated using [sebuf](https://github.com/SebastienMelki/sebuf) — generates HTTP handlers, type-safe clients, and OpenAPI docs.

```protobuf
service ExampleService {
  option (sebuf.http.service_config) = { base_path: "/api/v1" };
  rpc GetThing(GetThingRequest) returns (Thing) {
    option (sebuf.http.config) = { path: "/things/{id}" };
  }
}
```

Validation uses `buf.validate` annotations directly in proto files.

### Backend — DDD / Hexagonal Architecture

Each domain is a self-contained **module** under `backend/internal/modules/`:

```
modules/{domain}/
├── internal/
│   ├── domain/           # Entities, repository interfaces (ports), services
│   ├── application/      # Commands (writes) and queries (reads)
│   ├── infrastructure/   # Database implementations (adapters)
│   │   └── persistence/
│   │       └── sqlc/     # Generated type-safe SQL
│   └── interfaces/
│       └── http/         # HTTP handlers (sebuf-generated)
├── adapters/             # Cross-module adapters
├── sql/
│   ├── migrations/       # golang-migrate format
│   ├── queries/          # sqlc input
│   └── sqlc.yaml
└── module.go             # Wires everything together
```

**Key rules:**
- Domain layer has **zero external dependencies**
- Infrastructure implements domain interfaces (repository pattern)
- Modules communicate through **adapters**, never importing internals
- Each module owns its own migrations with separate tracking tables
- SQL queries are type-safe via [sqlc](https://sqlc.dev)

### Database

PostgreSQL via pgx/v5. For local dev and testing, use [testcontainers-go](https://golang.testcontainers.org/) to spin up disposable PostgreSQL instances. No embedded databases.

### Infrastructure

Empty Pulumi scaffold in Go. Add components as needed for your cloud provider.

---

## Development Commands

```bash
make generate       # Run all code generation (sqlc + buf)
make build          # Build everything
make test           # Run all tests
make lint           # Lint everything
make clean          # Clean build artifacts
```

### Backend
```bash
cd backend
make build          # Build server binary
make run            # Run server
make test           # Run tests (uses testcontainers, requires Docker)
make lint           # Run golangci-lint
make generate-sqlc  # Regenerate sqlc code
```

---

## Conventions

### Proto / Contracts
- One `.proto` file per message type, one `service.proto` per domain
- sebuf HTTP annotations for paths, `buf.validate` for validation
- Package pattern: `{project}/{domain}/v1`

### Go / Backend
- Module path: `github.com/{org}/{project}`
- Entry point: `cmd/server/main.go`
- Config via environment variables (godotenv in dev)
- Structured logging via slog
- Graceful shutdown on SIGINT/SIGTERM
- Health probes on port 8081

### Database / SQL
- Migrations: golang-migrate format, module-owned tracking tables
- All queries through sqlc — write SQL, get Go code

### Git
- Conventional commits
- Generated code IS committed
- `.env` files gitignored; `.env.example` committed

---

## Key Dependencies

| Tool | Purpose | Install |
|------|---------|---------|
| [buf](https://buf.build/docs/installation) | Proto management & code generation | `brew install bufbuild/buf/buf` |
| [sebuf](https://github.com/SebastienMelki/sebuf) | HTTP handler + client generation from proto | `go install github.com/SebastienMelki/sebuf/cmd/...@latest` |
| [sqlc](https://sqlc.dev) | Type-safe SQL code generation | `brew install sqlc` |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Database migrations | `brew install golang-migrate` |
| [golangci-lint](https://golangci-lint.run) | Go linting | `brew install golangci-lint` |
| [Pulumi](https://www.pulumi.com/docs/install/) | Infrastructure as code | `brew install pulumi` |

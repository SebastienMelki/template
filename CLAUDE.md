# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

This is a **GitHub template repository** for bootstrapping full-stack projects. When a developer creates a new repo from this template and opens Claude Code, you (Claude) should guide them through project setup by following the instructions below.

The template enforces a **contract-first, domain-driven** architecture with four pillars:

```
├── contracts/    # Protobuf API definitions (source of truth)
├── backend/      # Go service with DDD/hexagonal architecture
├── frontend/     # Web frontend (React by default, configurable)
└── infra/        # Pulumi infrastructure-as-code
```

---

## FIRST-TIME SETUP — Interactive Scaffolding

When a developer first opens this repo after creating it from the template, walk them through the following setup flow. Ask questions, then scaffold the project based on their answers.

### Phase 1: Project Identity

Ask the developer:
1. **Project name** — used for Go module path, package names, DNS, Pulumi stack names
2. **Organization/GitHub owner** — e.g. `github.com/acme` for Go module paths
3. **Brief project description** — for README and Pulumi project descriptions

### Phase 2: Contracts Configuration

Ask the developer:
1. **Domain modules** — what are the core business domains? (e.g. `user`, `order`, `product`)
   - For each domain, ask for 2-3 initial RPC methods they want to start with
2. **Shared types** — what identifiers will be shared across domains? (e.g. `UserID`, `OrderID`)
3. **Client targets** — which languages need generated clients?
   - Go (always, for the backend)
   - TypeScript (for web frontend)
   - Swift (for iOS)
   - Java (for Android)

### Phase 3: Backend Configuration

Ask the developer:
1. **Database** — PostgreSQL (default) or MySQL
2. **Database mode** — support in-memory dev mode? (recommended: yes)
3. **Additional integrations** — any of these needed from day one?
   - Slack notifications
   - LLM integration (OpenRouter/Claude)
   - Email (SMTP/SES)
   - File storage (S3-compatible)
4. **Auth strategy** — will be implemented as a module, but ask about approach:
   - Session-based
   - JWT
   - External provider (Clerk, Auth0)
   - None yet (add later)

### Phase 4: Frontend Configuration

Ask the developer:
1. **Framework** — React (default), or alternative:
   - React + Vite + TypeScript
   - Next.js + TypeScript
   - templ + HTMX (Go server-rendered, no separate frontend)
   - None (API-only project)
2. **UI library** — if React-based:
   - shadcn/ui + Tailwind (default)
   - Material UI
   - Ant Design
   - None / custom

### Phase 5: Infrastructure Configuration

Ask the developer:
1. **Cloud provider** — DigitalOcean (default) or AWS
2. **Compute** — Kubernetes (DOKS/EKS) or simple droplet/EC2
3. **Environments** — which stacks? (default: `dev` + `prod`)
4. **Domain** — do they have a domain? what is it?
5. **CI/CD** — GitHub Actions (default) with self-hosted runners?

### Phase 6: Scaffold

After collecting answers, scaffold the entire project. Replace all `{{PLACEHOLDER}}` values in template files. Create the directory structure, initial proto files, Go modules, frontend project, and Pulumi stacks. Then run through the **Post-Scaffold Checklist** below with the developer.

---

## POST-SCAFFOLD CHECKLIST

Walk the developer through these steps after scaffolding:

- [ ] `cd contracts && buf lint` — verify proto files are valid
- [ ] `make generate` — run code generation (proto → Go, TS, OpenAPI)
- [ ] `cd backend && make build` — verify backend compiles
- [ ] `cd backend && make test` — verify tests pass (in-memory DB mode)
- [ ] `cd frontend && npm install && npm run dev` — verify frontend starts (if applicable)
- [ ] `cd infra && pulumi preview --stack dev` — verify infra plan (if credentials configured)
- [ ] `git add -A && git commit -m "Initial scaffold from template"` — commit the scaffolded project
- [ ] Review generated README.md for accuracy

---

## Architecture

### Contract-First Development

All APIs are defined in Protocol Buffers under `contracts/`. The proto files are the **single source of truth**. Code is generated using [sebuf](https://github.com/SebastienMelki/sebuf) — a protobuf toolkit that generates HTTP handlers, type-safe clients, and OpenAPI docs.

**Flow:** Define proto → `buf generate` → generated Go HTTP handlers + TS clients + OpenAPI specs

sebuf annotations control HTTP behavior:
```protobuf
service ExampleService {
  option (sebuf.http.service_config) = { base_path: "/api/v1" };
  rpc GetThing(GetThingRequest) returns (Thing) {
    option (sebuf.http.config) = { path: "/things/{id}" };
  }
}
```

Validation uses `buf.validate` annotations directly in proto files — no separate validation layer needed.

### Backend — DDD / Hexagonal Architecture

Each business domain is a self-contained **module** under `backend/internal/modules/`. Modules follow this internal structure:

```
modules/{domain}/
├── internal/
│   ├── domain/           # Entities, value objects, repository interfaces (ports)
│   │   ├── entity/
│   │   ├── repository/   # Interfaces only — no implementations
│   │   └── service/      # Domain services (business logic)
│   ├── application/      # Use cases: commands (writes) and queries (reads)
│   │   ├── command/
│   │   ├── query/
│   │   └── mapper/       # DTO ↔ domain mappings
│   ├── infrastructure/   # Adapters: database implementations
│   │   └── persistence/
│   │       └── sqlc/     # Generated type-safe SQL
│   └── interfaces/
│       └── http/         # HTTP handlers (use generated sebuf code)
├── adapters/             # Cross-module adapters (ports exposed to other modules)
├── sql/
│   ├── migrations/       # Database schema (golang-migrate format)
│   ├── queries/          # SQL queries (sqlc input)
│   └── sqlc.yaml         # sqlc config for this module
└── module.go             # Module factory — wires everything together
```

**Key rules:**
- Domain layer has **zero external dependencies** — only stdlib and domain interfaces
- Infrastructure implements domain interfaces (repository pattern)
- Modules communicate through **adapters**, never by importing each other's internals
- Each module owns its own database migrations and schema
- SQL queries are type-safe via [sqlc](https://sqlc.dev) — write SQL, get Go code

### Frontend

When using React: standard Vite + TypeScript setup. Generated TypeScript clients from contracts are imported directly — no manual API layer needed.

When using templ + HTMX: templates live inside backend modules under `interfaces/templates/`, and the backend serves HTML directly.

### Infrastructure

Pulumi programs are written in Go (same language as backend). Each environment is a separate stack (`dev`, `prod`). Common components:
- VPC networking
- Managed database (PostgreSQL/MySQL)
- Kubernetes cluster or compute instance
- Container registry
- Ingress controller (Traefik) with Let's Encrypt TLS
- DNS (Cloudflare)
- CI/CD runner infrastructure

---

## Development Commands

### Root Level
```bash
make generate          # Run all code generation (proto + sqlc + templ)
make build             # Build everything
make test              # Run all tests
make lint              # Lint everything
make clean             # Clean generated files and build artifacts
```

### Contracts
```bash
cd contracts
buf lint               # Lint proto files
buf generate           # Generate code from proto definitions
buf breaking --against '.git#branch=main'  # Check for breaking changes
```

### Backend
```bash
cd backend
make build             # Build the server binary
make run               # Run with in-memory database
make test              # Run tests with race detection
make test-fast         # Run tests with caching (no race detection)
make lint              # Run golangci-lint
make migrate-up        # Apply database migrations (remote DB)
make migrate-down      # Rollback last migration
make fmt               # Format Go code
```

### Frontend (React)
```bash
cd frontend
npm install            # Install dependencies
npm run dev            # Start dev server
npm run build          # Production build
npm run lint           # Lint
npm run typecheck      # Type checking
```

### Infrastructure
```bash
cd infra
pulumi up --stack dev          # Deploy dev environment
pulumi preview --stack dev     # Preview changes
pulumi destroy --stack dev     # Tear down dev environment
make test                      # Run infra tests
```

---

## Conventions

### Proto / Contracts
- One `.proto` file per message type, one `service.proto` per domain service
- All services use sebuf HTTP annotations for path and header config
- Required fields must have `(buf.validate.field).required = true`
- Package pattern: `{project}/{domain}/v1`
- Service naming: `{Domain}Service` (e.g. `UserService`, `OrderService`)

### Go / Backend
- Module path: `github.com/{org}/{project}`
- Entry point: `cmd/server/main.go`
- Config via environment variables, loaded with godotenv in dev
- Structured logging via slog
- Database connection modes: `inmemory` (dev default), `remote` (production)
- Graceful shutdown on SIGINT/SIGTERM
- Health probes on separate port (default 8081)
- Each module registers its own HTTP routes and migrations

### Database / SQL
- Migrations use golang-migrate format: `{number}_{description}.up.sql` / `.down.sql`
- Each module has its own migration tracking table: `{module}_schema_migrations`
- All queries go through sqlc — no raw SQL in Go code
- Each module gets its own sqlc.yaml pointing to its queries/ and migrations/

### Git
- Conventional commits preferred
- Generated code IS committed (contracts generate once, commit the output)
- `.env` files are gitignored; `.env.example` is committed

---

## Key Dependencies

| Tool | Purpose | Install |
|------|---------|---------|
| [buf](https://buf.build/docs/installation) | Proto management & code generation | `brew install bufbuild/buf/buf` |
| [sebuf](https://github.com/SebastienMelki/sebuf) | HTTP handler + client generation from proto | `go install github.com/SebastienMelki/sebuf/cmd/...@latest` |
| [sqlc](https://sqlc.dev) | Type-safe SQL code generation | `brew install sqlc` |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Database migrations | `brew install golang-migrate` |
| [golangci-lint](https://golangci-lint.run) | Go linting | `brew install golangci-lint` |
| [templ](https://templ.guide) | Go HTML templates (if using HTMX frontend) | `go install github.com/a-h/templ/cmd/templ@latest` |
| [Pulumi](https://www.pulumi.com/docs/install/) | Infrastructure as code | `brew install pulumi` |

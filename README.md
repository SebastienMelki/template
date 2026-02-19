# {{PROJECT}}

{{PROJECT_DESCRIPTION}}

## Quick Start

Open this repo in [Claude Code](https://claude.ai/code) and say "Let's set it up." Claude will walk you through project configuration.

## Architecture

```
├── contracts/    # Protobuf API definitions (source of truth)
├── backend/      # Go service with DDD/hexagonal architecture
├── frontend/     # Web frontend
└── infra/        # Pulumi infrastructure-as-code
```

## Commands

```bash
make generate     # Generate code from proto + sqlc
make build        # Build everything
make test         # Run all tests
make lint         # Lint everything
```

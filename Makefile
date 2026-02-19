.PHONY: generate build test lint clean dev backend-run backend-test frontend-dev frontend-build infra-preview infra-up

# ============================================================================
# Code Generation
# ============================================================================

generate:
	buf generate
	$(MAKE) -C backend generate-sqlc

# ============================================================================
# Build
# ============================================================================

build:
	$(MAKE) -C backend build
	$(MAKE) -C frontend build

# ============================================================================
# Test
# ============================================================================

test:
	$(MAKE) -C backend test

# ============================================================================
# Lint
# ============================================================================

lint:
	$(MAKE) -C contracts lint
	$(MAKE) -C backend lint

# ============================================================================
# Clean
# ============================================================================

clean:
	$(MAKE) -C backend clean
	$(MAKE) -C frontend clean

# ============================================================================
# Development
# ============================================================================

dev:
	@echo "Starting backend and frontend..."
	$(MAKE) backend-run &
	$(MAKE) frontend-dev
	@wait

backend-run:
	$(MAKE) -C backend run

backend-test:
	$(MAKE) -C backend test

frontend-dev:
	$(MAKE) -C frontend dev

frontend-build:
	$(MAKE) -C frontend build

# ============================================================================
# Infrastructure
# ============================================================================

infra-preview:
	$(MAKE) -C infra preview

infra-up:
	$(MAKE) -C infra up

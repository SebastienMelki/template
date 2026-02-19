.PHONY: build run test test-fast lint fmt vet clean generate migrate-up migrate-down

BINARY_NAME=server
BUILD_DIR=bin

# ============================================================================
# Build & Run
# ============================================================================

build: generate
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

# ============================================================================
# Code Generation
# ============================================================================

generate: generate-sqlc

generate-sqlc:
	@for dir in internal/modules/*/sql; do \
		if [ -f "$$dir/sqlc.yaml" ]; then \
			echo "Generating sqlc for $$dir..."; \
			(cd "$$dir" && sqlc generate); \
		fi \
	done

# ============================================================================
# Test
# ============================================================================

test:
	go test -race -count=1 ./...

test-fast:
	go test ./...

# ============================================================================
# Code Quality
# ============================================================================

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# ============================================================================
# Database Migrations
# ============================================================================

migrate-up:
	@echo "Migrations are applied automatically on server startup"
	@echo "Set MIGRATION_ENABLED=true in .env"

# ============================================================================
# Clean
# ============================================================================

clean:
	rm -rf $(BUILD_DIR)

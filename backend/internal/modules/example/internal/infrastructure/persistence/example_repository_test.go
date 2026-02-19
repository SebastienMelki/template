package persistence_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/database"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/infrastructure/persistence"
)

const (
	testDBName = "testdb"
	testDBUser = "testuser"
	testDBPass = "testpass"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// setupTestDB spins up a PostgreSQL container via testcontainers and returns
// a database.DB wrapper plus a cleanup function.
func setupTestDB(t *testing.T) (*database.DB, func()) {
	t.Helper()

	if !dockerAvailable() {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pgxConn, err := pgx.Connect(ctx, connString)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Apply migration
	_, err = pgxConn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS examples (
			id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name        VARCHAR(255) NOT NULL,
			description TEXT,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to apply migration: %v", err)
	}

	db := &database.DB{Conn: pgxConn}

	cleanup := func() {
		pgxConn.Close(ctx)
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}

	return db, cleanup
}

func TestExampleRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := persistence.NewExampleRepository(db)
	ctx := context.Background()

	e := &entity.Example{
		ID:          uuid.New(),
		Name:        "Test Example",
		Description: "A test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, e)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify timestamps were set by DB
	if e.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestExampleRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := persistence.NewExampleRepository(db)
	ctx := context.Background()

	// Create an example first
	id := uuid.New()
	e := &entity.Example{
		ID:          id,
		Name:        "Get Test",
		Description: "Description for get test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := repo.Create(ctx, e); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get it back
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got == nil {
		t.Fatal("GetByID() returned nil")
	}
	if got.ID != id {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, id)
	}
	if got.Name != "Get Test" {
		t.Errorf("GetByID() Name = %v, want %v", got.Name, "Get Test")
	}
	if got.Description != "Description for get test" {
		t.Errorf("GetByID() Description = %v, want %v", got.Description, "Description for get test")
	}

	// Non-existent ID should return nil
	got, err = repo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID() for missing ID error = %v", err)
	}
	if got != nil {
		t.Error("GetByID() for missing ID should return nil")
	}
}

func TestExampleRepository_List(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := persistence.NewExampleRepository(db)
	ctx := context.Background()

	// Create several examples
	for i := 0; i < 5; i++ {
		e := &entity.Example{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("List Test %d", i),
			Description: fmt.Sprintf("Description %d", i),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := repo.Create(ctx, e); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	// List with pagination
	examples, err := repo.List(ctx, 3, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(examples) != 3 {
		t.Errorf("List(limit=3) returned %d items, want 3", len(examples))
	}

	// List second page
	examples, err = repo.List(ctx, 3, 3)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(examples) != 2 {
		t.Errorf("List(limit=3, offset=3) returned %d items, want 2", len(examples))
	}
}

func TestExampleRepository_Count(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := persistence.NewExampleRepository(db)
	ctx := context.Background()

	// Empty count
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Count() = %d, want 0", count)
	}

	// Create and count
	e := &entity.Example{
		ID:        uuid.New(),
		Name:      "Count Test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.Create(ctx, e); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	count, err = repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 1 {
		t.Errorf("Count() = %d, want 1", count)
	}
}

func TestExampleRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := persistence.NewExampleRepository(db)
	ctx := context.Background()

	// Create an example
	id := uuid.New()
	e := &entity.Example{
		ID:        id,
		Name:      "Delete Test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.Create(ctx, e); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Delete it
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's gone
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() after delete error = %v", err)
	}
	if got != nil {
		t.Error("GetByID() after delete should return nil")
	}
}

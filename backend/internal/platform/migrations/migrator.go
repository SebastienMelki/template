package migrations

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run applies all pending migrations from the given source directory.
// Each module calls this with its own migrations path and tracking table.
//
// The connString should be a PostgreSQL connection URL.
// The tableName is automatically prefixed: "{moduleName}_schema_migrations".
//
// Example:
//
//	migrations.Run(connString, "example", "file://internal/modules/example/sql/migrations")
func Run(connString string, moduleName string, sourceURL string) error {
	tableName := fmt.Sprintf("%s_schema_migrations", moduleName)

	// Append the migration table name to the connection string
	dbURL := fmt.Sprintf("%s&x-migrations-table=%s", connString, tableName)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator for %s: %w", moduleName, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations for %s: %w", moduleName, err)
	}

	version, dirty, _ := m.Version()
	slog.Info("migrations applied", "module", moduleName, "version", version, "dirty", dirty)

	return nil
}

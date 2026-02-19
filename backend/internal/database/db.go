package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX defines the interface that sqlc expects for PostgreSQL operations.
// Both DB and TX implement this interface, so repositories can work with either.
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// DB wraps pgx.Conn with instrumented logging.
type DB struct {
	*pgx.Conn
}

// Exec executes a SQL command and logs the operation.
func (db *DB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	slog.DebugContext(ctx, "executing SQL command", "query", sql)

	result, err := db.Conn.Exec(ctx, sql, arguments...)
	if err != nil {
		slog.ErrorContext(ctx, "SQL command failed", "query", sql, "error", err)
	}

	return result, err
}

// Query executes a SQL query that returns rows.
func (db *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.DebugContext(ctx, "executing SQL query", "query", sql)

	rows, err := db.Conn.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQL query failed", "query", sql, "error", err)
	}

	return rows, err
}

// QueryRow executes a SQL query that returns a single row.
func (db *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	slog.DebugContext(ctx, "executing SQL query (single row)", "query", sql)
	return db.Conn.QueryRow(ctx, sql, args...)
}

// TX wraps pgx.Tx with instrumented logging.
type TX struct {
	pgx.Tx
}

// Exec executes a SQL command within a transaction.
func (tx *TX) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	slog.DebugContext(ctx, "executing SQL command in transaction", "query", sql)

	result, err := tx.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		slog.ErrorContext(ctx, "SQL command failed in transaction", "query", sql, "error", err)
	}

	return result, err
}

// Query executes a SQL query that returns rows within a transaction.
func (tx *TX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.DebugContext(ctx, "executing SQL query in transaction", "query", sql)

	rows, err := tx.Tx.Query(ctx, sql, args...)
	if err != nil {
		slog.ErrorContext(ctx, "SQL query failed in transaction", "query", sql, "error", err)
	}

	return rows, err
}

// QueryRow executes a SQL query that returns a single row within a transaction.
func (tx *TX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	slog.DebugContext(ctx, "executing SQL query (single row) in transaction", "query", sql)
	return tx.Tx.QueryRow(ctx, sql, args...)
}

// BeginTx starts a new instrumented transaction.
func (db *DB) BeginTx(ctx context.Context, opts pgx.TxOptions) (*TX, error) {
	slog.DebugContext(ctx, "beginning database transaction")

	tx, err := db.Conn.BeginTx(ctx, opts)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TX{Tx: tx}, nil
}

// WithTransaction executes a function within a database transaction.
// The transaction is automatically committed if fn returns nil,
// or rolled back if it returns an error or panics.
func (db *DB) WithTransaction(ctx context.Context, fn func(DBTX) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.ErrorContext(ctx, "failed to rollback transaction after panic",
					"error", rollbackErr, "panic", r)
			}
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			slog.ErrorContext(ctx, "failed to rollback transaction",
				"rollback_error", rollbackErr, "original_error", err)
			return fmt.Errorf("transaction failed with rollback error: %w (original: %w)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close closes the underlying pgx connection.
func (db *DB) Close(ctx context.Context) error {
	err := db.Conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

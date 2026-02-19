// Package database provides database connection management for PostgreSQL.
// For local development and testing, use testcontainers-go to spin up a
// disposable PostgreSQL instance. In production, point to a remote database.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/jackc/pgx/v5"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/config"
)

// Connection represents a database connection with its associated resources.
type Connection struct {
	db         *DB
	connString string
}

// NewConnection creates a new database connection from the config.
func NewConnection(ctx context.Context, cfg config.DatabaseConfig) (*Connection, error) {
	connString := cfg.ConnectionString()

	slog.InfoContext(ctx, "connecting to PostgreSQL",
		"host", cfg.Host, "port", cfg.Port, "database", cfg.Name)

	pgxConn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx connection: %w", err)
	}

	if err := pgxConn.Ping(ctx); err != nil {
		if closeErr := pgxConn.Close(ctx); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close pgx connection after ping error", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.InfoContext(ctx, "connected to PostgreSQL")

	return &Connection{
		db:         &DB{Conn: pgxConn},
		connString: connString,
	}, nil
}

// NewConnectionFromString creates a connection from a raw connection string.
// Useful for testcontainers where the connection string is provided dynamically.
func NewConnectionFromString(ctx context.Context, connString string) (*Connection, error) {
	pgxConn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx connection: %w", err)
	}

	if err := pgxConn.Ping(ctx); err != nil {
		if closeErr := pgxConn.Close(ctx); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close pgx connection after ping error", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{
		db:         &DB{Conn: pgxConn},
		connString: connString,
	}, nil
}

// DB returns the instrumented DB wrapper for application use.
func (c *Connection) DB() *DB {
	return c.db
}

// GetConnectionString returns the connection string (useful for migration tools).
func (c *Connection) GetConnectionString() string {
	return c.connString
}

// Close closes the database connection.
func (c *Connection) Close(ctx context.Context) error {
	slog.InfoContext(ctx, "closing database connection")

	if c.db != nil {
		if err := c.db.Close(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to close database connection", "error", err)
			return err
		}
	}

	return nil
}

// BuildConnectionString builds a PostgreSQL connection string from individual parts.
// Useful for testcontainers or other dynamic connection setups.
func BuildConnectionString(host string, port int, user, password, dbName, sslMode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		dbName,
		sslMode)
}

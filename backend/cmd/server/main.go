package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/config"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/database"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/health"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/logging"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize logging
	logging.Setup(cfg.LogLevel, cfg.DebugMode)

	slog.Info("starting server", "port", cfg.Port)

	// Initialize database
	conn, err := database.NewConnection(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Initialize modules
	exampleModule := example.New(conn.DB())

	// Run migrations
	if cfg.Database.MigrationEnabled {
		if err := exampleModule.Migrate(conn.GetConnectionString()); err != nil {
			slog.Error("failed to run migrations", "module", "example", "error", err)
			os.Exit(1)
		}
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	exampleModule.RegisterRoutes(mux)

	// Main server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Probe server (health checks)
	probeMux := http.NewServeMux()
	health.RegisterRoutes(probeMux)
	probeServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ProbePort),
		Handler:           probeMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start servers
	go func() {
		slog.Info("probe server listening", "port", cfg.ProbePort)
		if err := probeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("probe server failed", "error", err)
		}
	}()

	go func() {
		slog.Info("main server listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("main server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
	if err := probeServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("probe server shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}

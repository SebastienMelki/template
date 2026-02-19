package example

import (
	"net/http"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/database"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/service"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/infrastructure/persistence"
	handler "github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/interfaces/http"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/platform/migrations"
)

// Module is the example module's public API.
// It wires together domain, infrastructure, and interface layers.
type Module struct {
	handler *handler.Handler
	service *service.ExampleService
}

// New creates a new example module with all dependencies wired.
func New(db *database.DB) *Module {
	repo := persistence.NewExampleRepository(db)
	svc := service.NewExampleService(repo)
	h := handler.NewHandler(svc)

	return &Module{
		handler: h,
		service: svc,
	}
}

// RegisterRoutes registers this module's HTTP routes on the given mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	m.handler.RegisterRoutes(mux)
}

// Migrate runs this module's database migrations.
func (m *Module) Migrate(connString string) error {
	return migrations.Run(connString, "example",
		"file://internal/modules/example/sql/migrations")
}

// Service returns the underlying service for adapter creation.
func (m *Module) Service() *service.ExampleService {
	return m.service
}

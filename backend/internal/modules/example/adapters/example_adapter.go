package adapters

import (
	"context"

	"github.com/google/uuid"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/service"
)

// ExampleAdapter exposes example module functionality to other modules.
// This is the cross-module boundary — other modules depend on this interface,
// not on the internal domain directly.
type ExampleAdapter struct {
	service *service.ExampleService
}

// NewExampleAdapter creates a new adapter.
func NewExampleAdapter(svc *service.ExampleService) *ExampleAdapter {
	return &ExampleAdapter{service: svc}
}

// GetByID retrieves an example by ID, exposed for cross-module use.
func (a *ExampleAdapter) GetByID(ctx context.Context, id uuid.UUID) (*entity.Example, error) {
	return a.service.GetByID(ctx, id)
}

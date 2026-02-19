package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/service"
)

// GetExampleQuery retrieves a single example by ID.
type GetExampleQuery struct {
	ID uuid.UUID
}

// GetExampleHandler handles the GetExample query.
type GetExampleHandler struct {
	service *service.ExampleService
}

// NewGetExampleHandler creates a new handler.
func NewGetExampleHandler(svc *service.ExampleService) *GetExampleHandler {
	return &GetExampleHandler{service: svc}
}

// Handle executes the query.
func (h *GetExampleHandler) Handle(ctx context.Context, q GetExampleQuery) (*entity.Example, error) {
	return h.service.GetByID(ctx, q.ID)
}

// ListExamplesQuery retrieves a paginated list of examples.
type ListExamplesQuery struct {
	Limit  int32
	Offset int32
}

// ListExamplesHandler handles listing examples.
type ListExamplesHandler struct {
	service *service.ExampleService
}

// NewListExamplesHandler creates a new handler.
func NewListExamplesHandler(svc *service.ExampleService) *ListExamplesHandler {
	return &ListExamplesHandler{service: svc}
}

// Handle executes the query.
func (h *ListExamplesHandler) Handle(ctx context.Context, q ListExamplesQuery) ([]*entity.Example, int64, error) {
	return h.service.List(ctx, q.Limit, q.Offset)
}

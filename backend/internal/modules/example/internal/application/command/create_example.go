package command

import (
	"context"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/service"
)

// CreateExampleCommand represents a request to create a new example.
type CreateExampleCommand struct {
	Name        string
	Description string
}

// CreateExampleHandler handles the CreateExample command.
type CreateExampleHandler struct {
	service *service.ExampleService
}

// NewCreateExampleHandler creates a new handler.
func NewCreateExampleHandler(svc *service.ExampleService) *CreateExampleHandler {
	return &CreateExampleHandler{service: svc}
}

// Handle executes the command and returns the created entity.
func (h *CreateExampleHandler) Handle(ctx context.Context, cmd CreateExampleCommand) (*entity.Example, error) {
	return h.service.Create(ctx, cmd.Name, cmd.Description)
}

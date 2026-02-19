package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
)

// ExampleRepository defines the persistence interface for the example domain.
// This is a PORT — it lives in the domain layer and has no external dependencies.
type ExampleRepository interface {
	Create(ctx context.Context, e *entity.Example) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Example, error)
	List(ctx context.Context, limit, offset int32) ([]*entity.Example, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

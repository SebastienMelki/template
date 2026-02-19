package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/repository"
)

// ExampleService contains the business logic for the example domain.
type ExampleService struct {
	repo repository.ExampleRepository
}

// NewExampleService creates a new service with the given repository.
func NewExampleService(repo repository.ExampleRepository) *ExampleService {
	return &ExampleService{repo: repo}
}

// Create creates a new example entity.
func (s *ExampleService) Create(ctx context.Context, name, description string) (*entity.Example, error) {
	now := time.Now()
	e := &entity.Example{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	return e, nil
}

// GetByID retrieves an example by its ID.
func (s *ExampleService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Example, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves examples with pagination.
func (s *ExampleService) List(ctx context.Context, limit, offset int32) ([]*entity.Example, int64, error) {
	examples, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list examples: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count examples: %w", err)
	}

	return examples, total, nil
}

// Delete removes an example by its ID.
func (s *ExampleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

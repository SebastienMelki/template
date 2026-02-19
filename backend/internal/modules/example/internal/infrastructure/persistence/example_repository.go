package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/database"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/domain/entity"
	"github.com/{{ORG}}/{{PROJECT}}/backend/internal/modules/example/internal/infrastructure/persistence/sqlc"
)

// ExampleRepository is the PostgreSQL implementation of the domain repository.
// It uses sqlc-generated queries for type-safe database access.
type ExampleRepository struct {
	queries *sqlc.Queries
}

// NewExampleRepository creates a new repository backed by the given DBTX.
func NewExampleRepository(db database.DBTX) *ExampleRepository {
	return &ExampleRepository{
		queries: sqlc.New(db),
	}
}

// Create persists a new example entity.
func (r *ExampleRepository) Create(ctx context.Context, e *entity.Example) error {
	var desc pgtype.Text
	if e.Description != "" {
		desc = pgtype.Text{String: e.Description, Valid: true}
	}

	result, err := r.queries.CreateExample(ctx, sqlc.CreateExampleParams{
		ID:          uuidToPgtype(e.ID),
		Name:        e.Name,
		Description: desc,
	})
	if err != nil {
		return fmt.Errorf("failed to insert example: %w", err)
	}

	e.CreatedAt = result.CreatedAt.Time
	e.UpdatedAt = result.UpdatedAt.Time
	return nil
}

// GetByID retrieves an example by its ID.
func (r *ExampleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Example, error) {
	row, err := r.queries.GetExampleByID(ctx, uuidToPgtype(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get example: %w", err)
	}
	return sqlcToEntity(row), nil
}

// List retrieves examples with pagination.
func (r *ExampleRepository) List(ctx context.Context, limit, offset int32) ([]*entity.Example, error) {
	rows, err := r.queries.ListExamples(ctx, sqlc.ListExamplesParams{
		ResultLimit:  limit,
		ResultOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list examples: %w", err)
	}

	examples := make([]*entity.Example, 0, len(rows))
	for _, row := range rows {
		examples = append(examples, sqlcToEntity(row))
	}
	return examples, nil
}

// Count returns the total number of examples.
func (r *ExampleRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountExamples(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count examples: %w", err)
	}
	return count, nil
}

// Delete removes an example by its ID.
func (r *ExampleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteExample(ctx, uuidToPgtype(id)); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}
	return nil
}

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func sqlcToEntity(row sqlc.Example) *entity.Example {
	desc := ""
	if row.Description.Valid {
		desc = row.Description.String
	}
	return &entity.Example{
		ID:          uuid.UUID(row.ID.Bytes),
		Name:        row.Name,
		Description: desc,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

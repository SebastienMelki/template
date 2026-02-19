package entity

import (
	"time"

	"github.com/google/uuid"
)

// Example is the core domain entity.
type Example struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

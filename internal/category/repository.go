package category

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	List(ctx context.Context) ([]Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
}

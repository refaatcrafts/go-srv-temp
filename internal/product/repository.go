package product

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, categoryID *uuid.UUID) ([]Product, error)
}

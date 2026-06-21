package postgres

import (
	"context"
	"fmt"
	"strings"

	"go-srv-temp/internal/product"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) product.Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *product.Product) error {
	query := `INSERT INTO products (id, name, description, price, currency, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`
	return r.db.QueryRowContext(ctx, query, p.ID, p.Name, p.Description, p.Price, p.Currency, p.CategoryID).Scan(&p.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	var p product.Product
	query := `SELECT id, name, description, price, currency, category_id, created_at FROM products WHERE id = $1`
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context, categoryID *uuid.UUID) ([]product.Product, error) {
	var args []any
	var conditions []string

	if categoryID != nil {
		args = append(args, *categoryID)
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}

	query := `SELECT id, name, description, price, currency, category_id, created_at FROM products`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC"

	var pp []product.Product
	err := r.db.SelectContext(ctx, &pp, query, args...)
	return pp, err
}

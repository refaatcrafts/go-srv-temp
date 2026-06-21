package postgres

import (
	"context"
	"errors"
	"fmt"

	"go-srv-temp/internal/category"
	"go-srv-temp/internal/httperr"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) category.Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *category.Category) error {
	query := `INSERT INTO categories (id, name, slug) VALUES ($1, $2, $3) RETURNING created_at`
	err := r.db.QueryRowContext(ctx, query, c.ID, c.Name, c.Slug).Scan(&c.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w", httperr.ErrDuplicateKey)
		}
		return err
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	var c category.Category
	query := `SELECT id, name, slug, created_at FROM categories WHERE id = $1`
	err := r.db.GetContext(ctx, &c, query, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(ctx context.Context) ([]category.Category, error) {
	var cc []category.Category
	query := `SELECT id, name, slug, created_at FROM categories ORDER BY name`
	err := r.db.SelectContext(ctx, &cc, query)
	return cc, err
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (*category.Category, error) {
	var c category.Category
	query := `SELECT id, name, slug, created_at FROM categories WHERE slug = $1`
	err := r.db.GetContext(ctx, &c, query, slug)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

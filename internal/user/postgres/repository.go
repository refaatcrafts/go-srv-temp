package postgres

import (
	"context"
	"errors"
	"fmt"

	"go-srv-temp/internal/httperr"
	"go-srv-temp/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) user.Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING created_at`
	err := r.db.QueryRowContext(ctx, query, u.ID, u.Email, u.PasswordHash).Scan(&u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w", httperr.ErrDuplicateKey)
		}
		return err
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

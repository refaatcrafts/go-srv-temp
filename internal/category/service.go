package category

import (
	"context"
	"database/sql"
	"errors"

	"go-srv-temp/internal/httperr"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, httperr.New(422, "name and slug are required")
	}

	c := &Category{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		if errors.Is(err, httperr.ErrDuplicateKey) {
			return nil, httperr.ErrConflict
		}
		return nil, err
	}

	return c, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if 		errors.Is(err, sql.ErrNoRows) {
			return nil, httperr.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) List(ctx context.Context) ([]Category, error) {
	return s.repo.List(ctx)
}

func (s *Service) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if 		errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

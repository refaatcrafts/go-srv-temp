package product

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"go-srv-temp/internal/httperr"

	"github.com/google/uuid"
)

type CategoryService interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type Service struct {
	repo            Repository
	categoryService CategoryService
}

func NewService(repo Repository, cs CategoryService) *Service {
	return &Service{repo: repo, categoryService: cs}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	if req.Name == "" {
		return nil, httperr.New(422, "name is required")
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, httperr.New(422, "invalid category_id")
	}

	exists, err := s.categoryService.Exists(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, httperr.New(422, "category not found")
	}

	if req.Price < 0 {
		return nil, httperr.New(422, "price must be non-negative")
	}

	priceCents := int64(math.Round(req.Price * 100))
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	p := &Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       priceCents,
		Currency:    currency,
		CategoryID:  categoryID,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if 		errors.Is(err, sql.ErrNoRows) {
			return nil, httperr.ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) List(ctx context.Context, categoryID *uuid.UUID) ([]Product, error) {
	pp, err := s.repo.List(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if pp == nil {
		pp = []Product{}
	}
	return pp, nil
}

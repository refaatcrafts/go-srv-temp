package product_test

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"testing"

	"go-srv-temp/internal/httperr"
	"go-srv-temp/internal/product"

	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc func(ctx context.Context, p *product.Product) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*product.Product, error)
	listFunc    func(ctx context.Context, categoryID *uuid.UUID) ([]product.Product, error)
}

func (m *mockRepo) Create(ctx context.Context, p *product.Product) error {
	return m.createFunc(ctx, p)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) List(ctx context.Context, categoryID *uuid.UUID) ([]product.Product, error) {
	return m.listFunc(ctx, categoryID)
}

type mockCatSvc struct {
	existsFunc func(ctx context.Context, id uuid.UUID) (bool, error)
}

func (m *mockCatSvc) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return m.existsFunc(ctx, id)
}

func TestService_Create(t *testing.T) {
	catID := uuid.New()

	tests := []struct {
		name    string
		req     product.CreateProductRequest
		catSvc  *mockCatSvc
		wantErr bool
		errCode int
	}{
		{
			name: "empty name",
			req: product.CreateProductRequest{
				Name:       "",
				Price:      10.0,
				CategoryID: catID.String(),
			},
			wantErr: true,
			errCode: 422,
		},
		{
			name: "invalid category_id",
			req: product.CreateProductRequest{
				Name:       "test",
				Price:      10.0,
				CategoryID: "not-a-uuid",
			},
			wantErr: true,
			errCode: 422,
		},
		{
			name: "category not found",
			req: product.CreateProductRequest{
				Name:       "test",
				Price:      10.0,
				CategoryID: catID.String(),
			},
			catSvc: &mockCatSvc{
				existsFunc: func(ctx context.Context, id uuid.UUID) (bool, error) {
					return false, nil
				},
			},
			wantErr: true,
			errCode: 422,
		},
		{
			name: "negative price",
			req: product.CreateProductRequest{
				Name:       "test",
				Price:      -1,
				CategoryID: catID.String(),
			},
			catSvc: &mockCatSvc{
				existsFunc: func(ctx context.Context, id uuid.UUID) (bool, error) {
					return true, nil
				},
			},
			wantErr: true,
			errCode: 422,
		},
		{
			name: "success with USD default",
			req: product.CreateProductRequest{
				Name:       "Test Product",
				Price:      19.99,
				CategoryID: catID.String(),
			},
			catSvc: &mockCatSvc{
				existsFunc: func(ctx context.Context, id uuid.UUID) (bool, error) {
					return true, nil
				},
			},
			wantErr: false,
		},
		{
			name: "success with EUR",
			req: product.CreateProductRequest{
				Name:       "Test Product",
				Price:      9.50,
				Currency:   "EUR",
				CategoryID: catID.String(),
			},
			catSvc: &mockCatSvc{
				existsFunc: func(ctx context.Context, id uuid.UUID) (bool, error) {
					return true, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				createFunc: func(ctx context.Context, p *product.Product) error {
					return nil
				},
			}
			svc := product.NewService(repo, tt.catSvc)

			p, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var e *httperr.Error
				if errors.As(err, &e) && e.Code != tt.errCode {
					t.Errorf("expected status %d, got %d", tt.errCode, e.Code)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.Name != tt.req.Name {
				t.Errorf("expected name %q, got %q", tt.req.Name, p.Name)
			}
			expectedCents := int64(math.Round(tt.req.Price * 100))
			if p.Price != expectedCents {
				t.Errorf("expected price %d cents, got %d", expectedCents, p.Price)
			}
			if p.CategoryID != catID {
				t.Errorf("expected category_id %v, got %v", catID, p.CategoryID)
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	prodID := uuid.New()

	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr bool
		errType error
	}{
		{
			name: "not found",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*product.Product, error) {
					return nil, sql.ErrNoRows
				},
			},
			wantErr: true,
			errType: httperr.ErrNotFound,
		},
		{
			name: "success",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*product.Product, error) {
					return &product.Product{ID: id, Name: "test"}, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := product.NewService(tt.repo, nil)
			p, err := svc.GetByID(context.Background(), prodID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tt.errType) {
					t.Errorf("expected %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.ID != prodID {
				t.Errorf("expected id %v, got %v", prodID, p.ID)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockRepo
		categoryID *uuid.UUID
		wantNil    bool
	}{
		{
			name: "nil list returns empty slice",
			repo: &mockRepo{
				listFunc: func(ctx context.Context, categoryID *uuid.UUID) ([]product.Product, error) {
					return nil, nil
				},
			},
			categoryID: nil,
		},
		{
			name: "returns products",
			repo: &mockRepo{
				listFunc: func(ctx context.Context, categoryID *uuid.UUID) ([]product.Product, error) {
					return []product.Product{{Name: "p1"}}, nil
				},
			},
			categoryID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := product.NewService(tt.repo, nil)
			pp, err := svc.List(context.Background(), tt.categoryID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if pp == nil {
				t.Error("expected non-nil slice, got nil")
			}
		})
	}
}

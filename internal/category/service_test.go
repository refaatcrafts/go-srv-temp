package category_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-srv-temp/internal/category"
	"go-srv-temp/internal/httperr"

	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc  func(ctx context.Context, c *category.Category) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*category.Category, error)
	listFunc    func(ctx context.Context) ([]category.Category, error)
	getBySlugFunc func(ctx context.Context, slug string) (*category.Category, error)
}

func (m *mockRepo) Create(ctx context.Context, c *category.Category) error {
	return m.createFunc(ctx, c)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) List(ctx context.Context) ([]category.Category, error) {
	return m.listFunc(ctx)
}
func (m *mockRepo) GetBySlug(ctx context.Context, slug string) (*category.Category, error) {
	return m.getBySlugFunc(ctx, slug)
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name    string
		req     category.CreateCategoryRequest
		repo    *mockRepo
		wantErr bool
		errType error
	}{
		{
			name: "empty name",
			req: category.CreateCategoryRequest{
				Name: "",
				Slug: "test",
			},
			wantErr: true,
		},
		{
			name: "empty slug",
			req: category.CreateCategoryRequest{
				Name: "test",
				Slug: "",
			},
			wantErr: true,
		},
		{
			name: "duplicate slug",
			req: category.CreateCategoryRequest{
				Name: "test",
				Slug: "test",
			},
			repo: &mockRepo{
				createFunc: func(ctx context.Context, c *category.Category) error {
					return httperr.ErrDuplicateKey
				},
			},
			wantErr: true,
			errType: httperr.ErrConflict,
		},
		{
			name: "success",
			req: category.CreateCategoryRequest{
				Name: "Electronics",
				Slug: "electronics",
			},
			repo: &mockRepo{
				createFunc: func(ctx context.Context, c *category.Category) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := category.NewService(tt.repo)
			c, err := svc.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.Name != tt.req.Name {
				t.Errorf("expected name %q, got %q", tt.req.Name, c.Name)
			}
			if c.Slug != tt.req.Slug {
				t.Errorf("expected slug %q, got %q", tt.req.Slug, c.Slug)
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	catID := uuid.New()

	tests := []struct {
		name    string
		repo    *mockRepo
		wantErr bool
		errType error
	}{
		{
			name: "not found",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*category.Category, error) {
					return nil, sql.ErrNoRows
				},
			},
			wantErr: true,
			errType: httperr.ErrNotFound,
		},
		{
			name: "success",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*category.Category, error) {
					return &category.Category{ID: id, Name: "test"}, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := category.NewService(tt.repo)
			c, err := svc.GetByID(context.Background(), catID)

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
			if c.ID != catID {
				t.Errorf("expected id %v, got %v", catID, c.ID)
			}
		})
	}
}

func TestService_Exists(t *testing.T) {
	catID := uuid.New()

	tests := []struct {
		name    string
		repo    *mockRepo
		want    bool
		wantErr bool
	}{
		{
			name: "exists",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*category.Category, error) {
					return &category.Category{ID: id}, nil
				},
			},
			want: true,
		},
		{
			name: "not found",
			repo: &mockRepo{
				getByIDFunc: func(ctx context.Context, id uuid.UUID) (*category.Category, error) {
					return nil, sql.ErrNoRows
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := category.NewService(tt.repo)
			exists, err := svc.Exists(context.Background(), catID)
			if tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if exists != tt.want {
				t.Errorf("expected %v, got %v", tt.want, exists)
			}
		})
	}
}

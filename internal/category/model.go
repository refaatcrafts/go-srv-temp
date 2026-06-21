package category

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       int64     `db:"price" json:"price"`
	Currency    string    `db:"currency" json:"currency"`
	CategoryID  uuid.UUID `db:"category_id" json:"category_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	CategoryID  string  `json:"category_id"`
}

package category

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/categories", h.Create)
	r.Get("/categories", h.List)
	r.Get("/categories/{id}", h.GetByID)
}

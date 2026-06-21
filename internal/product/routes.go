package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler, authMw func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authMw)
		r.Post("/products", h.Create)
	})

	r.Get("/products", h.List)
	r.Get("/products/{id}", h.GetByID)
}

package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler, authMw func(http.Handler) http.Handler) {
	r.Post("/users/signup", h.Signup)
	r.Post("/users/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMw)
		r.Get("/users/me", h.Me)
	})
}

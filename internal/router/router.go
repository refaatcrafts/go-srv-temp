package router

import (
	"go-srv-temp/internal/category"
	"go-srv-temp/internal/middleware"
	"go-srv-temp/internal/product"
	"go-srv-temp/internal/user"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(
	logger *slog.Logger,
	jwtSecret string,
	categoryHandler *category.Handler,
	userHandler *user.Handler,
	productHandler *product.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recoverer(logger))

	authMw := middleware.Auth(jwtSecret)

	r.Route("/api/v1", func(r chi.Router) {
		category.RegisterRoutes(r, categoryHandler)
		user.RegisterRoutes(r, userHandler, authMw)
		product.RegisterRoutes(r, productHandler, authMw)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

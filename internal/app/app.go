package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-srv-temp/internal/category"
	categorypg "go-srv-temp/internal/category/postgres"
	"go-srv-temp/internal/config"
	"go-srv-temp/internal/product"
	productpg "go-srv-temp/internal/product/postgres"
	"go-srv-temp/internal/router"
	"go-srv-temp/internal/user"
	userpg "go-srv-temp/internal/user/postgres"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *sqlx.DB
	server *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db, err := sqlx.Open("pgx", cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	logger.Info("database connected")

	catRepo := categorypg.NewRepository(db)
	catSvc := category.NewService(catRepo)
	catHandler := category.NewHandler(catSvc)

	jwtSvc := user.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	userRepo := userpg.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtSvc)
	userHandler := user.NewHandler(userSvc)

	prodRepo := productpg.NewRepository(db)
	prodSvc := product.NewService(prodRepo, catSvc)
	prodHandler := product.NewHandler(prodSvc)

	handler := router.New(logger, cfg.JWT.Secret, catHandler, userHandler, prodHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		a.logger.Info("server starting", "port", a.cfg.App.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("server error", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	a.logger.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	a.logger.Info("server stopped")
	return nil
}

package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-srv-temp/internal/httperr"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo  Repository
	jwt   *JWTService
}

func NewService(repo Repository, jwt *JWTService) *Service {
	return &Service{repo: repo, jwt: jwt}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		return nil, httperr.New(422, "email and password are required")
	}
	if len(req.Password) < 8 {
		return nil, httperr.New(422, "password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(ctx, u); err != nil {
		if errors.Is(err, httperr.ErrDuplicateKey) {
			return nil, httperr.ErrConflict
		}
		return nil, err
	}

	token, err := s.jwt.GenerateToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *u}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		return nil, httperr.New(422, "email and password are required")
	}

	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if 		errors.Is(err, sql.ErrNoRows) {
			return nil, httperr.New(401, "invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, httperr.New(401, "invalid email or password")
	}

	token, err := s.jwt.GenerateToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *u}, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if 		errors.Is(err, sql.ErrNoRows) {
			return nil, httperr.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

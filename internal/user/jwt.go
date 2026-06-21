package user

import (
	"time"

	"go-srv-temp/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret    string
	expiresIn time.Duration
}

func NewJWTService(secret string, expiresIn time.Duration) *JWTService {
	return &JWTService{secret: secret, expiresIn: expiresIn}
}

func (s *JWTService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

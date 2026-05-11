package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// AuthUserRepo is the minimal user-repository contract needed by AuthService.
// Defined here (where it's consumed) so Login can be unit-tested against fakes.
// *repository.MongoTenantRepo satisfies this interface structurally.
type AuthUserRepo interface {
	GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, id string, t time.Time) error
}

type AuthService struct {
	userRepo  AuthUserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo AuthUserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type TokenClaims struct {
	UserID   string      `json:"sub"`
	TenantID string      `json:"tid"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expiresAt"`
	TenantID  string      `json:"tenantId"`
	Role      domain.Role `json:"role"`
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	if !user.Active {
		return nil, domain.ErrForbidden
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := TokenClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now())

	return &LoginResponse{
		Token:     signed,
		ExpiresAt: expiresAt,
		TenantID:  user.TenantID,
		Role:      user.Role,
	}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

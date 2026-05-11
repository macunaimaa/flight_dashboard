package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
)

// newTestService builds an AuthService without a repository.
// Safe for ValidateToken / HashPassword tests since neither touches the repo.
func newTestService(secret string) *AuthService {
	return &AuthService{jwtSecret: []byte(secret)}
}

func signClaims(t *testing.T, secret string, claims TokenClaims, method jwt.SigningMethod) string {
	t.Helper()
	tok := jwt.NewWithClaims(method, claims)
	signed, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign claims: %v", err)
	}
	return signed
}

func TestHashPassword_RoundTrip(t *testing.T) {
	hash, err := HashPassword("hunter2")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "hunter2" {
		t.Fatal("hash should differ from plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("hunter2")); err != nil {
		t.Fatalf("correct password rejected: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrong")); err == nil {
		t.Fatal("incorrect password accepted")
	}
}

func TestHashPassword_DifferentSaltsProduceDifferentHashes(t *testing.T) {
	h1, err := HashPassword("samepw")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	h2, err := HashPassword("samepw")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if h1 == h2 {
		t.Fatal("expected distinct salts to produce distinct hashes")
	}
}

func TestValidateToken_HappyPath(t *testing.T) {
	svc := newTestService("test-secret")
	claims := TokenClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     domain.Role("admin"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user-1",
		},
	}
	signed := signClaims(t, "test-secret", claims, jwt.SigningMethodHS256)

	got, err := svc.ValidateToken(signed)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if got.UserID != "user-1" || got.TenantID != "tenant-1" || got.Role != "admin" {
		t.Fatalf("claims mismatch: %+v", got)
	}
}

func TestValidateToken_RejectsExpired(t *testing.T) {
	svc := newTestService("test-secret")
	claims := TokenClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     domain.Role("viewer"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	signed := signClaims(t, "test-secret", claims, jwt.SigningMethodHS256)

	_, err := svc.ValidateToken(signed)
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for expired token, got %v", err)
	}
}

func TestValidateToken_RejectsWrongSecret(t *testing.T) {
	svc := newTestService("real-secret")
	claims := TokenClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     domain.Role("admin"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	signed := signClaims(t, "other-secret", claims, jwt.SigningMethodHS256)

	_, err := svc.ValidateToken(signed)
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for wrong-secret token, got %v", err)
	}
}

func TestValidateToken_RejectsMalformed(t *testing.T) {
	svc := newTestService("test-secret")
	_, err := svc.ValidateToken("not-a-jwt")
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for malformed token, got %v", err)
	}
}

func TestValidateToken_RejectsNoneAlg(t *testing.T) {
	svc := newTestService("test-secret")
	// Build a token with alg=none — should never validate even if the signature is empty.
	claims := TokenClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     domain.Role("admin"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}

	_, err = svc.ValidateToken(signed)
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for alg=none token, got %v", err)
	}
}

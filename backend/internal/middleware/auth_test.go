package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

const testSecret = "test-jwt-secret"

type spy struct {
	called   bool
	tenantID string
	userID   string
	role     domain.Role
}

func newAuthSvc() *service.AuthService {
	return service.NewAuthService(nil, testSecret)
}

func signTestToken(t *testing.T, claims service.TokenClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return signed
}

func spyHandler(s *spy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.called = true
		if v, ok := r.Context().Value(service.CtxKeyTenantID).(string); ok {
			s.tenantID = v
		}
		if v, ok := r.Context().Value(service.CtxKeyUserID).(string); ok {
			s.userID = v
		}
		if v, ok := r.Context().Value(service.CtxKeyRole).(domain.Role); ok {
			s.role = v
		}
		w.WriteHeader(http.StatusOK)
	})
}

func runMiddleware(t *testing.T, header string) (*httptest.ResponseRecorder, *spy) {
	t.Helper()
	s := &spy{}
	h := Auth(newAuthSvc())(spyHandler(s))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec, s
}

func TestAuth_MissingHeader(t *testing.T) {
	rec, s := runMiddleware(t, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if s.called {
		t.Fatal("next handler should not be called when header is missing")
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	rec, s := runMiddleware(t, "Token abc.def.ghi") // no "Bearer " prefix
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if s.called {
		t.Fatal("next handler should not be called for non-Bearer format")
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	rec, s := runMiddleware(t, "Bearer not-a-valid-jwt")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if s.called {
		t.Fatal("next handler should not be called for invalid token")
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	token := signTestToken(t, service.TokenClaims{
		UserID:   "u1",
		TenantID: "t1",
		Role:     domain.Role("admin"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	})

	rec, s := runMiddleware(t, "Bearer "+token)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if s.called {
		t.Fatal("next handler should not be called for expired token")
	}
}

func TestAuth_ValidTokenInjectsContext(t *testing.T) {
	token := signTestToken(t, service.TokenClaims{
		UserID:   "user-42",
		TenantID: "tenant-42",
		Role:     domain.Role("admin"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	rec, s := runMiddleware(t, "Bearer "+token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !s.called {
		t.Fatal("next handler should be called for valid token")
	}
	if s.tenantID != "tenant-42" {
		t.Errorf("ctx tenantID = %q, want tenant-42", s.tenantID)
	}
	if s.userID != "user-42" {
		t.Errorf("ctx userID = %q, want user-42", s.userID)
	}
	if s.role != domain.Role("admin") {
		t.Errorf("ctx role = %q, want admin", s.role)
	}
}

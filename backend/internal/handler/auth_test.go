package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

func newAuthHandler(t *testing.T, u *fakeAuthUserRepo) *AuthHandler {
	t.Helper()
	svc := service.NewAuthService(u, "test-secret")
	return NewAuthHandler(svc)
}

func mustBcrypt(t *testing.T, pw string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(h)
}

func postJSON(h http.HandlerFunc, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func TestLogin_MalformedJSON(t *testing.T) {
	repo := &fakeAuthUserRepo{}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, "{not-json")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	repo := &fakeAuthUserRepo{}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, `{"email":"","password":"","tenantId":""}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := &fakeAuthUserRepo{userErr: domain.ErrNotFound}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, `{"email":"x@y.com","password":"pw","tenantId":"t1"}`)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	repo := &fakeAuthUserRepo{
		user: &domain.User{
			ID:           "u1",
			TenantID:     "t1",
			Email:        "x@y.com",
			PasswordHash: mustBcrypt(t, "correct"),
			Role:         domain.RoleAdmin,
			Active:       false,
		},
	}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, `{"email":"x@y.com","password":"correct","tenantId":"t1"}`)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &fakeAuthUserRepo{
		user: &domain.User{
			ID:           "u1",
			TenantID:     "t1",
			Email:        "x@y.com",
			PasswordHash: mustBcrypt(t, "correct"),
			Role:         domain.RoleAdmin,
			Active:       true,
		},
	}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, `{"email":"x@y.com","password":"WRONG","tenantId":"t1"}`)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestLogin_HappyPath(t *testing.T) {
	repo := &fakeAuthUserRepo{
		user: &domain.User{
			ID:           "u1",
			TenantID:     "t1",
			Email:        "x@y.com",
			PasswordHash: mustBcrypt(t, "correct"),
			Role:         domain.RoleAdmin,
			Active:       true,
		},
	}
	h := newAuthHandler(t, repo)

	rec := postJSON(h.Login, `{"email":"x@y.com","password":"correct","tenantId":"t1"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}

	var resp APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("body not JSON: %v", err)
	}
	if resp.Error != nil {
		t.Fatalf("Error = %v, want nil", *resp.Error)
	}
	d, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data not a map: %T", resp.Data)
	}
	tok, _ := d["token"].(string)
	if tok == "" || !strings.Contains(tok, ".") {
		t.Errorf("token field looks wrong: %q", tok)
	}
	if d["tenantId"] != "t1" {
		t.Errorf("tenantId = %v, want t1", d["tenantId"])
	}
	if d["role"] != string(domain.RoleAdmin) {
		t.Errorf("role = %v, want %s", d["role"], domain.RoleAdmin)
	}

	// Fire-and-forget UpdateLastLogin runs in a goroutine; allow it to settle.
	// We don't assert it deterministically — the goroutine could race the test exit.
	_ = repo
}

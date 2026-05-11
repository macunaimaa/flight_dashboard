package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body not JSON: %v (body=%s)", err, rec.Body.String())
	}
	if body.Error != nil {
		t.Errorf("error = %v, want nil", *body.Error)
	}
	dataMap, ok := body.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data not a map, got %T", body.Data)
	}
	if dataMap["status"] != "ok" {
		t.Errorf("data.status = %v, want ok", dataMap["status"])
	}
}

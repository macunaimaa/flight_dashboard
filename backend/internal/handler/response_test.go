package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestWriteData(t *testing.T) {
	rec := httptest.NewRecorder()
	writeData(rec, map[string]int{"count": 7})

	if rec.Code != 200 {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var resp APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if resp.Error != nil {
		t.Errorf("Error = %v, want nil", *resp.Error)
	}
	if resp.Meta == nil {
		t.Fatal("Meta is nil")
	}
	if resp.Meta.Timestamp.IsZero() {
		t.Error("Meta.Timestamp is zero")
	}
	d, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data not a map: %T", resp.Data)
	}
	if d["count"].(float64) != 7 {
		t.Errorf("Data.count = %v, want 7", d["count"])
	}
}

func TestWriteDataWithMeta(t *testing.T) {
	rec := httptest.NewRecorder()
	writeDataWithMeta(rec, []int{1, 2, 3}, 42, 10, 5)

	var resp APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if resp.Meta.Total != 42 {
		t.Errorf("Total = %d, want 42", resp.Meta.Total)
	}
	if resp.Meta.Limit != 10 {
		t.Errorf("Limit = %d, want 10", resp.Meta.Limit)
	}
	if resp.Meta.Offset != 5 {
		t.Errorf("Offset = %d, want 5", resp.Meta.Offset)
	}
	arr, ok := resp.Data.([]interface{})
	if !ok || len(arr) != 3 {
		t.Errorf("Data = %v, want length-3 array", resp.Data)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, 418, "I'm a teapot")

	if rec.Code != 418 {
		t.Fatalf("status = %d, want 418", rec.Code)
	}
	var resp APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if resp.Error == nil {
		t.Fatal("Error is nil, want set")
	}
	if *resp.Error != "I'm a teapot" {
		t.Errorf("Error = %q, want I'm a teapot", *resp.Error)
	}
	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
}

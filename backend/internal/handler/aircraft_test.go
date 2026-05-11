package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

const testTenant = "tenant-test"

func newAircraftDeps() (*fakeAircraftRepo, *fakeCacheRepo, *fakeAuditRepo, *AircraftHandler) {
	repo := &fakeAircraftRepo{}
	cache := &fakeCacheRepo{}
	audit := &fakeAuditRepo{}
	aircraftSvc := service.NewAircraftService(repo, cache)
	auditSvc := service.NewAuditService(audit)
	return repo, cache, audit, NewAircraftHandler(aircraftSvc, auditSvc)
}

func withTenant(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), service.CtxKeyTenantID, testTenant)
	ctx = context.WithValue(ctx, service.CtxKeyUserID, "user-test")
	return req.WithContext(ctx)
}

func TestList_DefaultFilters(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	repo.listStates = []domain.AircraftLiveState{{ICAO24: "abc123"}}
	repo.listTotal = 1

	req := withTenant(httptest.NewRequest(http.MethodGet, "/api/v1/aircraft", nil))
	rec := httptest.NewRecorder()

	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	if repo.listGotFilter.Limit != 100 {
		t.Errorf("default limit = %d, want 100", repo.listGotFilter.Limit)
	}
	if repo.listGotFilter.Offset != 0 {
		t.Errorf("default offset = %d, want 0", repo.listGotFilter.Offset)
	}
	if repo.listGotFilter.Callsign != nil {
		t.Errorf("Callsign should be nil, got %v", repo.listGotFilter.Callsign)
	}
}

func TestList_WiresQueryParams(t *testing.T) {
	repo, _, _, h := newAircraftDeps()

	req := withTenant(httptest.NewRequest(http.MethodGet,
		"/api/v1/aircraft?limit=25&offset=50&callsign=UAL&country=US&on_ground=true", nil))
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	f := repo.listGotFilter
	if f.Limit != 25 || f.Offset != 50 {
		t.Errorf("limit/offset = %d/%d, want 25/50", f.Limit, f.Offset)
	}
	if f.Callsign == nil || *f.Callsign != "UAL" {
		t.Errorf("Callsign = %v, want UAL", f.Callsign)
	}
	if f.Country == nil || *f.Country != "US" {
		t.Errorf("Country = %v, want US", f.Country)
	}
	if f.OnGround == nil || *f.OnGround != true {
		t.Errorf("OnGround = %v, want true", f.OnGround)
	}
}

func TestList_RepoErrorReturns500(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	repo.listErr = errAircraft("boom")

	req := withTenant(httptest.NewRequest(http.MethodGet, "/api/v1/aircraft", nil))
	rec := httptest.NewRecorder()
	h.List(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestGetByICAO_HappyPath(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	repo.byICAO = &domain.AircraftLiveState{ICAO24: "abc123", Callsign: "UAL999"}

	r := chi.NewRouter()
	r.Get("/aircraft/{icao24}", h.GetByICAO)

	req := withTenant(httptest.NewRequest(http.MethodGet, "/aircraft/abc123", nil))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	var resp APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	d, _ := resp.Data.(map[string]interface{})
	if d["icao24"] != "abc123" {
		t.Errorf("icao24 = %v, want abc123", d["icao24"])
	}

	if len(repo.byICAOCalls) != 1 {
		t.Fatalf("byICAOCalls = %d, want 1", len(repo.byICAOCalls))
	}
	if repo.byICAOCalls[0].tenant != testTenant {
		t.Errorf("tenant = %q, want %q", repo.byICAOCalls[0].tenant, testTenant)
	}
}

func TestGetByICAO_NotFound(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	repo.byICAOErr = domain.ErrNotFound

	r := chi.NewRouter()
	r.Get("/aircraft/{icao24}", h.GetByICAO)

	req := withTenant(httptest.NewRequest(http.MethodGet, "/aircraft/missing", nil))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestGetInBBox_ParsesParams(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	repo.bboxStates = []domain.AircraftLiveState{{ICAO24: "x"}, {ICAO24: "y"}}

	req := withTenant(httptest.NewRequest(http.MethodGet,
		"/api/v1/aircraft/bbox?lamin=10&lomin=20&lamax=30&lomax=40&limit=50", nil))
	rec := httptest.NewRecorder()
	h.GetInBBox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	if len(repo.bboxCalls) != 1 {
		t.Fatalf("bboxCalls = %d, want 1", len(repo.bboxCalls))
	}
	c := repo.bboxCalls[0]
	if c.bbox.LaMin != 10 || c.bbox.LoMin != 20 || c.bbox.LaMax != 30 || c.bbox.LoMax != 40 {
		t.Errorf("bbox = %+v, want (10,20,30,40)", c.bbox)
	}
	if c.limit != 50 {
		t.Errorf("limit = %d, want 50", c.limit)
	}
}

func TestGetInBBox_DefaultsWhenParamsMissing(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	req := withTenant(httptest.NewRequest(http.MethodGet, "/api/v1/aircraft/bbox", nil))
	rec := httptest.NewRecorder()
	h.GetInBBox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if len(repo.bboxCalls) != 1 {
		t.Fatalf("bboxCalls = %d, want 1", len(repo.bboxCalls))
	}
	c := repo.bboxCalls[0]
	if c.bbox.LaMin != -90 || c.bbox.LoMin != -180 || c.bbox.LaMax != 90 || c.bbox.LoMax != 180 {
		t.Errorf("default bbox = %+v, want full earth", c.bbox)
	}
	if c.limit != 1000 {
		t.Errorf("default limit = %d, want 1000", c.limit)
	}
}

func TestGetTrack_ParsesParams(t *testing.T) {
	repo, _, _, h := newAircraftDeps()
	from := "2026-01-01T00:00:00Z"
	to := "2026-01-02T00:00:00Z"

	r := chi.NewRouter()
	r.Get("/aircraft/{icao24}/track", h.GetTrack)

	req := withTenant(httptest.NewRequest(http.MethodGet,
		"/aircraft/abc123/track?from="+from+"&to="+to, nil))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	if len(repo.trackCalls) != 1 {
		t.Fatalf("trackCalls = %d, want 1", len(repo.trackCalls))
	}
	got := repo.trackCalls[0]
	if got.icao != "abc123" {
		t.Errorf("icao = %q, want abc123", got.icao)
	}
	wantFrom, _ := time.Parse(time.RFC3339, from)
	wantTo, _ := time.Parse(time.RFC3339, to)
	if !got.from.Equal(wantFrom) {
		t.Errorf("from = %v, want %v", got.from, wantFrom)
	}
	if !got.to.Equal(wantTo) {
		t.Errorf("to = %v, want %v", got.to, wantTo)
	}
}

func TestParseHelpers(t *testing.T) {
	if parseIntParam("", 7) != 7 {
		t.Error("parseIntParam empty default")
	}
	if parseIntParam("42", 0) != 42 {
		t.Error("parseIntParam valid")
	}
	if parseIntParam("nope", 99) != 99 {
		t.Error("parseIntParam invalid falls back to default")
	}
	if parseFloatParam("", 1.5) != 1.5 {
		t.Error("parseFloatParam empty default")
	}
	if parseFloatParam("3.14", 0) != 3.14 {
		t.Error("parseFloatParam valid")
	}
	if parseFloatParam("nope", 9.9) != 9.9 {
		t.Error("parseFloatParam invalid falls back to default")
	}
	def := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if !parseTimeParam("", def).Equal(def) {
		t.Error("parseTimeParam empty default")
	}
	if parseTimeParam("not-a-time", def).Equal(def) == false {
		t.Error("parseTimeParam invalid default")
	}
	got := parseTimeParam("2026-05-11T12:00:00Z", def)
	want, _ := time.Parse(time.RFC3339, "2026-05-11T12:00:00Z")
	if !got.Equal(want) {
		t.Errorf("parseTimeParam valid = %v, want %v", got, want)
	}
}

// errAircraft is a sentinel error type for repo test failures.
type errAircraft string

func (e errAircraft) Error() string { return string(e) }

package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

type AircraftHandler struct {
	aircraftSvc *service.AircraftService
	auditSvc    *service.AuditService
}

func NewAircraftHandler(aircraftSvc *service.AircraftService, auditSvc *service.AuditService) *AircraftHandler {
	return &AircraftHandler{
		aircraftSvc: aircraftSvc,
		auditSvc:    auditSvc,
	}
}

func (h *AircraftHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := service.TenantIDFromCtx(r.Context())
	q := r.URL.Query()

	filter := domain.AircraftFilter{
		Limit:  parseIntParam(q.Get("limit"), 100),
		Offset: parseIntParam(q.Get("offset"), 0),
	}
	if cs := q.Get("callsign"); cs != "" {
		filter.Callsign = &cs
	}
	if c := q.Get("country"); c != "" {
		filter.Country = &c
	}
	if og := q.Get("on_ground"); og != "" {
		val := og == "true"
		filter.OnGround = &val
	}

	states, total, err := h.aircraftSvc.GetLiveStates(r.Context(), tenantID, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch aircraft")
		return
	}

	h.auditSvc.Log(r.Context(), tenantID, service.UserIDFromCtx(r.Context()),
		"aircraft.list", "aircraft_live_state", filter)

	writeDataWithMeta(w, states, total, filter.Limit, filter.Offset)
}

func (h *AircraftHandler) GetByICAO(w http.ResponseWriter, r *http.Request) {
	tenantID := service.TenantIDFromCtx(r.Context())
	icao24 := chi.URLParam(r, "icao24")

	state, err := h.aircraftSvc.GetByICAO(r.Context(), tenantID, icao24)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "aircraft not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch aircraft")
		return
	}

	writeData(w, state)
}

func (h *AircraftHandler) GetInBBox(w http.ResponseWriter, r *http.Request) {
	tenantID := service.TenantIDFromCtx(r.Context())
	q := r.URL.Query()

	bbox := domain.BBox{
		LaMin: parseFloatParam(q.Get("lamin"), -90),
		LoMin: parseFloatParam(q.Get("lomin"), -180),
		LaMax: parseFloatParam(q.Get("lamax"), 90),
		LoMax: parseFloatParam(q.Get("lomax"), 180),
	}
	limit := parseIntParam(q.Get("limit"), 1000)

	states, err := h.aircraftSvc.GetInBBox(r.Context(), tenantID, bbox, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch aircraft")
		return
	}

	writeDataWithMeta(w, states, int64(len(states)), limit, 0)
}

func (h *AircraftHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	tenantID := service.TenantIDFromCtx(r.Context())
	icao24 := chi.URLParam(r, "icao24")
	q := r.URL.Query()

	from := parseTimeParam(q.Get("from"), time.Now().Add(-1*time.Hour))
	to := parseTimeParam(q.Get("to"), time.Now())

	points, err := h.aircraftSvc.GetTrackPoints(r.Context(), tenantID, icao24, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch track")
		return
	}

	writeDataWithMeta(w, points, int64(len(points)), len(points), 0)
}

func parseIntParam(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func parseFloatParam(s string, def float64) float64 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}

func parseTimeParam(s string, def time.Time) time.Time {
	if s == "" {
		return def
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return def
	}
	return t
}

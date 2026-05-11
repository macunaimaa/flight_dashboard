package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Data  interface{} `json:"data"`
	Meta  *APIMeta    `json:"meta,omitempty"`
	Error *string     `json:"error,omitempty"`
}

type APIMeta struct {
	Total     int64     `json:"total,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeData(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{
		Data: data,
		Meta: &APIMeta{Timestamp: time.Now()},
	})
}

func writeDataWithMeta(w http.ResponseWriter, data interface{}, total int64, limit, offset int) {
	writeJSON(w, http.StatusOK, APIResponse{
		Data: data,
		Meta: &APIMeta{
			Total:     total,
			Limit:     limit,
			Offset:    offset,
			Timestamp: time.Now(),
		},
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIResponse{Error: &message})
}

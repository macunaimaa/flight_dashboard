package handler

import (
	"net/http"

	"github.com/macunaimaa/dashboard/backend/internal/service"
)

type NewsHandler struct {
	newsSvc *service.NewsService
}

func NewNewsHandler(newsSvc *service.NewsService) *NewsHandler {
	return &NewsHandler{newsSvc: newsSvc}
}

func (h *NewsHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "missing query parameter 'q'")
		return
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = "en"
	}

	result, err := h.newsSvc.Search(r.Context(), query, lang)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeData(w, result)
}

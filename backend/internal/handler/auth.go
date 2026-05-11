package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/macunaimaa/dashboard/backend/internal/domain"
	"github.com/macunaimaa/dashboard/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.TenantID == "" {
		writeError(w, http.StatusBadRequest, "email, password and tenantId are required")
		return
	}

	resp, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "account disabled")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeData(w, resp)
}

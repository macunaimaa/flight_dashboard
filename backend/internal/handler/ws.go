package handler

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/macunaimaa/dashboard/backend/internal/service"
	"github.com/macunaimaa/dashboard/backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS handled by middleware
	},
}

type WSHandler struct {
	hub     *ws.Hub
	authSvc *service.AuthService
	logger  *slog.Logger
}

func NewWSHandler(hub *ws.Hub, authSvc *service.AuthService, logger *slog.Logger) *WSHandler {
	return &WSHandler{
		hub:     hub,
		authSvc: authSvc,
		logger:  logger,
	}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Auth via query param (WebSocket can't set custom headers)
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := h.authSvc.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("ws upgrade failed", "error", err)
		return
	}

	client := ws.NewClient(h.hub, conn, claims.TenantID, claims.UserID, h.logger)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

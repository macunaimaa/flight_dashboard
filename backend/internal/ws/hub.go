package ws

import (
	"log/slog"
	"sync"
)

type Hub struct {
	tenantClients map[string]map[*Client]bool
	register      chan *Client
	unregister    chan *Client
	broadcast     chan tenantMessage
	logger        *slog.Logger
	mu            sync.RWMutex
}

type tenantMessage struct {
	TenantID string
	Data     []byte
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		tenantClients: make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan tenantMessage, 256),
		logger:        logger.With("component", "ws_hub"),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.tenantClients[client.TenantID] == nil {
				h.tenantClients[client.TenantID] = make(map[*Client]bool)
			}
			h.tenantClients[client.TenantID][client] = true
			count := len(h.tenantClients[client.TenantID])
			h.mu.Unlock()

			h.logger.Info("client connected",
				"tenant_id", client.TenantID,
				"user_id", client.UserID,
				"clients_count", count,
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.tenantClients[client.TenantID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
				}
				if len(clients) == 0 {
					delete(h.tenantClients, client.TenantID)
				}
			}
			h.mu.Unlock()

			h.logger.Info("client disconnected",
				"tenant_id", client.TenantID,
				"user_id", client.UserID,
			)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.tenantClients[msg.TenantID]
			for client := range clients {
				select {
				case client.send <- msg.Data:
				default:
					// Buffer full — disconnect slow client
					close(client.send)
					delete(clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) BroadcastToTenant(tenantID string, data []byte) {
	h.broadcast <- tenantMessage{TenantID: tenantID, Data: data}
}

func (h *Hub) ClientCount(tenantID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.tenantClients[tenantID])
}

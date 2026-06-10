package ws

import (
	"encoding/json"
	"sync"

	gows "github.com/gofiber/websocket/v2"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type client struct {
	conn     *gows.Conn
	tenantID string
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*client]bool
}

var Default = &Hub{clients: make(map[*client]bool)}

func (h *Hub) Register(conn *gows.Conn, tenantID string) {
	c := &client{conn: conn, tenantID: tenantID}
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.mu.Lock()
			delete(h.clients, c)
			h.mu.Unlock()
			return
		}
	}
}

func (h *Hub) Broadcast(tenantID string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.tenantID == tenantID {
			_ = c.conn.WriteMessage(gows.TextMessage, data)
		}
	}
}

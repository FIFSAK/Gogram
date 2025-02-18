package ws

import (
	"fmt"
	"github.com/FIFSAK/Gogram/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[int64]*websocket.Conn
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int64]*websocket.Conn),
	}
}

func (h *Hub) RegisterConnection(userID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[userID] = conn
}

func (h *Hub) UnregisterConnection(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, userID)
}

func (h *Hub) SendMessage(userID int64, message models.Message) error {
	h.mu.RLock()
	conn, ok := h.connections[userID]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("user %d not connected", userID)
	}

	if err := conn.WriteJSON(message); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

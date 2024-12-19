package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn *websocket.Conn
	User string // User ID
}

type Hub struct {
	Connections map[string]*Connection
	mu          sync.Mutex
}

var Instance = Hub{
	Connections: make(map[string]*Connection),
}

func (h *Hub) AddConnection(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.Connections[userID] = &Connection{Conn: conn, User: userID}
	h.mu.Unlock()
}

func (h *Hub) RemoveConnection(userID string) {
	h.mu.Lock()
	delete(h.Connections, userID)
	h.mu.Unlock()
}

func (h *Hub) BroadcastMessage(userID string, message interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for id, conn := range h.Connections {
		if id != userID {
			conn.Conn.WriteJSON(message)
		}
	}
}

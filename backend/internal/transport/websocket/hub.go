package websocket

import (
	"log/slog"
	"sync"
)

// https://echo.labstack.com/cookbook/websocket/

type Hub struct {
	// concurrent safe map: key=UserID, value=*websocket.Conn
	clients sync.Map
}

func NewHub(log *slog.Logger) *Hub {
	return &Hub{}
}

// func (h *Hub) Connect(userID int64, conn *websocket.Conn) {
// 	h.clients.Store(userID, conn)
// }

// func (h *Hub) Disconnect(userID int64) {
// 	if conn, ok := h.clients.LoadAndDelete(userID); ok {
// 		conn.(*websocket.Conn).Close()
// 	}
// }

func (h *Hub) IsOnline(userID int64) bool {
	_, online := h.clients.Load(userID)
	return online
}

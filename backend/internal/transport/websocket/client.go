package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       int64  `json:"id"` // Identifies the authenticated owner -> UserID
	Username string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Conn     *websocket.Conn
	Message  chan *Message // TODO: it shouldn't be message, maybe []byte?
	// Hub      *ConnHub      // ???
}

type Message struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id,omitempty"`
	Username string `json:"username"`
	Content  string `json:"content"`
	// System    bool   `json:"system"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// ReadMessage
func (cl *Client) ReadMessage(hub *ConnHub) {
	defer func() {
		hub.Unregister <- cl
		cl.Conn.Close()
	}()

	for {
		_, m, err := cl.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:   string(m),
			ChatID:    cl.ChatID,
			Username:  cl.Username,
			SenderID:  cl.ID,
			Timestamp: time.Now(),
		}

		hub.Broadcast <- msg
	}
}

// WriteMessage
func (cl *Client) WriteMessage() {
	defer func() {
		cl.Conn.Close()
	}()

	for {
		message, ok := <-cl.Message
		if !ok {
			return
		}

		cl.Conn.WriteJSON(message)
	}
}

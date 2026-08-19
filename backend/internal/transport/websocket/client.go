package websocket

import (
	"encoding/json"
	"log/slog"
	"match-me-api/internal/logger"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second    // max time allowed to write a message frame to the peer
	pongWait       = 60 * time.Second    // max time allowed to wait for the next pong response from the peer
	pingPeriod     = (pongWait * 9) / 10 // 54 seconds (must be strictly less than pongWait).
	maxMessageSize = 4096                // size (in bytes) of an incoming message
	sendBufferSize = 256                 // outgoing message buffer per client
)

// Client represents a single active WebSocket connection for an authenticated user.
type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
	hub    *Hub
	log    *slog.Logger
}

// NewClient instantiates a new Client session.
func NewClient(userID int64, conn *websocket.Conn, hub *Hub, logger *slog.Logger) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, sendBufferSize),
		hub:    hub,
		log:    logger,
	}
}

// ClientInboundMessage wraps an incoming event with the sender's client session.
type ClientInboundMessage struct {
	Client *Client
	Event  InboundEvent
}

// ReadPump reads incoming WebSocket messages from the connection and dispatches them to the Hub.
// There is at most one reader per connection running in a dedicated goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Debug("websocket connection closed unexpectedly", "user_id", c.UserID, logger.Err(err))
			}
			break
		}

		var event InboundEvent
		if err := json.Unmarshal(rawMessage, &event); err != nil {
			c.SendError("invalid JSON event payload")
			continue
		}

		c.hub.Inbound <- &ClientInboundMessage{
			Client: c,
			Event:  event,
		}
	}
}

// WritePump writes queued messages from the Send channel to the WebSocket connection.
// There is at most one writer per connection running in a dedicated goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Hub closed the channel; close the WebSocket connection gracefully
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				return
			}

			// Flush any additional queued messages in the buffer in a single frame
			n := len(c.Send)
			for range n {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					return
				}
				if _, err := w.Write(<-c.Send); err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendEvent serializes an OutboundEvent to JSON and enqueues it to the client's Send channel.
func (c *Client) SendEvent(event OutboundEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		c.log.Error("failed to marshal outbound websocket event", logger.Err(err))
		return
	}

	select {
	case c.Send <- data:
	default:
		// Client buffer is full; drop connection to prevent blocking the hub
		c.log.Warn("client buffer full, dropping message", "user_id", c.UserID)
	}
}

// SendError is a helper to send an error event back to the client.
func (c *Client) SendError(msg string) {
	c.SendEvent(OutboundEvent{
		Type:    EventError,
		Payload: ErrorPayload{Message: msg},
	})
}


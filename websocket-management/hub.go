package websocketManagement

import (
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 1024),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			if oldClient, exists := h.clients[c.username]; exists {
				h.closeClient(oldClient, "Connection closed because another session connection opened with this username")
			}
			h.clients[c.username] = c

		case c := <-h.unregister:
			current, ok := h.clients[c.username]
			if ok && current == c {
				delete(h.clients, c.username)
				h.closeClient(c, "disconnected")
			}

		case msg := <-h.broadcast:
			for _, c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// slow client → drop
					delete(h.clients, c.username)
					close(c.send)
				}
			}
		}
	}
}

func (h *Hub) closeClient(c *Client, reason string) {
	msg := websocket.FormatCloseMessage(
		websocket.CloseNormalClosure,
		reason,
	)

	_ = c.conn.WriteControl(
		websocket.CloseMessage,
		msg,
		time.Now().Add(time.Second),
	)

	c.conn.Close()
	close(c.send)
}

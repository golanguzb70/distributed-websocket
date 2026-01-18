package websocketManagement

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan MessageSent
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan MessageSent, 1024),
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
			fmt.Printf("New connection is established for user %s\n", c.username)
		case c := <-h.unregister:
			current, ok := h.clients[c.username]
			if ok && current == c {
				delete(h.clients, c.username)
				h.closeClient(c, "disconnected")
				fmt.Printf("Connection of user %s is closed\n", c.username)
			}

		case msg := <-h.broadcast:
			msgRec := MessageReceive{
				From:    msg.From,
				Message: msg.Message,
			}

			bt, err := json.MarshalIndent(msgRec, "", "\t")
			if err != nil {
				fmt.Println(err)
			}

			for _, username := range msg.To {
				c, ok := h.clients[username]
				if ok {
					select {
					case c.send <- bt:
					default:
						// slow client → drop
						delete(h.clients, c.username)
						close(c.send)
					}
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

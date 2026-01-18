package websocketManagement

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan []byte
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- msg
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

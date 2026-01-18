package websocketManagement

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan []byte
}

type MessageSent struct {
	To      []string `json:"to"`
	From    string   `json:"from"`
	Message string   `json:"message"`
}

type MessageReceive struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type MessageRedis struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		msgObj := MessageSent{}
		err = json.Unmarshal(msg, &msgObj)
		msgObj.From = c.username
		h.broadcast <- msgObj
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

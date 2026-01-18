package websocketManagement

import (
	"distributed-websocket/redis"
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
	redisMsg   chan []byte
	redis      *redis.RedisConn
}

func NewHub(redisConn *redis.RedisConn) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan MessageSent, 1024),
		redisMsg:   make(chan []byte, 1024),
		redis:      redisConn,
	}
}

func (h *Hub) Run() {
	go h.redis.SubscribeAndWriteToChann(h.redisMsg)

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
				fmt.Println("Error while marshaling", err)
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
				} else {
					msgRedis := MessageRedis{
						To:      username,
						From:    msg.From,
						Message: msg.Message,
					}
					bt, err := json.MarshalIndent(msgRedis, "", "\t")
					if err != nil {
						fmt.Println("Error while marshaling redis message", err)
						continue
					}
					err = h.redis.Publish(bt)
					if err != nil {
						fmt.Println("Erorr while publishing message", err)
					}
				}
			}
		case msg := <-h.redisMsg:
			msgObj := MessageRedis{}
			err := json.Unmarshal(msg, &msgObj)
			if err != nil {
				fmt.Println("Couldn't unmarshal redis msg", err, string(msg))
				continue
			}
			c, ok := h.clients[msgObj.To]
			if ok {
				msgRec := MessageReceive{
					From:    msgObj.From,
					Message: msgObj.Message,
				}

				bt, err := json.MarshalIndent(msgRec, "", "\t")
				if err != nil {
					fmt.Println("error while marshaling msg res", err)
					continue
				}
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

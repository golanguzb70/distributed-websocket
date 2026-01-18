package main

import (
	websocketManagement "distributed-websocket/websocket-management"
	"fmt"
	"net/http"
)

func main() {
	hub := websocketManagement.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketManagement.ServeWS(hub, w, r)
	})

	fmt.Println("Websocket is listening on :8080 port...")
	http.ListenAndServe(":8080", nil)
}

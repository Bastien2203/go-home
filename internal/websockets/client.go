package websockets

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // TODO remove true for production
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan *Message
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("JSON Error: %v", err)
			continue
		}

		switch msg.Action {
		case "subscribe":
			c.hub.Subscribe(c, msg.Topic)
			log.Printf("Client subscribed to : %s", msg.Topic)

		case "publish":
			c.hub.broadcast <- &msg
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		msg, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteJSON(msg); err != nil {
			return
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan *Message, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

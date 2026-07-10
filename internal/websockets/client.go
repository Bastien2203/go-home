package websockets

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Bastien2203/go-home/shared/config"
	"github.com/Bastien2203/go-home/shared/middlewares"
	"github.com/gorilla/websocket"
)

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

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, allowedOrigins []string, appEnv config.AppEnv) {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = true
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			if appEnv == config.Dev && middlewares.IsLocalhostOrigin(origin) {
				return true
			}
			return originSet[origin]
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan *Message, 256)}

	go client.writePump()
	go client.readPump()
}

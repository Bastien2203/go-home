package websockets

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	topics map[Topic]map[*Client]bool

	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		topics:     make(map[Topic]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 100),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:

			_ = client

		case client := <-h.unregister:
			h.mu.Lock()

			for topic := range h.topics {
				if _, ok := h.topics[topic][client]; ok {
					delete(h.topics[topic], client)
					close(client.send)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()

			if clients, ok := h.topics[msg.Topic]; ok {
				for client := range clients {
					select {
					case client.send <- msg:
					default:
						close(client.send)
						delete(h.topics[msg.Topic], client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(topic Topic, payload any) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &Message{
		Action:  "broadcast",
		Topic:   topic,
		Message: jsonPayload,
	}

	h.broadcast <- msg

	return nil
}

func (h *Hub) Subscribe(client *Client, topic Topic) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.topics[topic]; !ok {
		h.topics[topic] = make(map[*Client]bool)
	}
	h.topics[topic][client] = true
}

package websockets

import "encoding/json"

type Message struct {
	Action  string          `json:"action"`
	Topic   Topic           `json:"topic"`
	Message json.RawMessage `json:"message"`
}

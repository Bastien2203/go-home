package events

import "sync"

type EventType string

const (
	RawDataReceived EventType = "raw_data_received"
)

type Event struct {
	Type    EventType
	Payload any
}

type EventHandler func(event Event)

type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if handlers, exists := eb.handlers[event.Type]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

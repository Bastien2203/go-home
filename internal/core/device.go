package core

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"time"
)

type Event struct {
	DeviceAddr string
	Type       string
	Payload    map[string]any
	Time       time.Time
}

type Device interface {
	Addr() string
	Name() string
	Type() DeviceType
	State() map[string]any
	SetState(k string, v any)
	Subscribe(ch chan<- Event)
	Unsubscribe(ch chan<- Event)
	Start(ctx context.Context) error
	ToJson() map[string]any
}

type DeviceType string

const (
	TemperatureSensorType DeviceType = "temperature_sensor"
	HumiditySensorType    DeviceType = "humidity_sensor"
)

var DeviceTypes = map[DeviceType]string{
	TemperatureSensorType: "Temperature Sensor",
	HumiditySensorType:    "Humidity Sensor",
}

// BaseDevice provides thread-safe state + pub/sub
type BaseDevice struct {
	addr  string
	name  string
	kind  DeviceType
	mu    sync.RWMutex
	state map[string]any

	subsMu  sync.RWMutex
	subs    map[chan<- Event]struct{}
	parser  Parser
	running bool
}

func NewBaseDevice(addr string, name string, kind DeviceType, parser Parser) *BaseDevice {
	return &BaseDevice{
		addr:    addr,
		name:    name,
		kind:    kind,
		state:   make(map[string]any),
		subs:    make(map[chan<- Event]struct{}),
		parser:  parser,
		running: false,
	}
}

func (b *BaseDevice) ToJson() map[string]any {
	return map[string]any{
		"addr":        b.addr,
		"name":        b.name,
		"type":        b.Type(),
		"parser_type": b.parser.Name(),
		"running":     b.running,
	}
}

func (b *BaseDevice) Addr() string     { return b.addr }
func (b *BaseDevice) Name() string     { return b.name }
func (b *BaseDevice) Type() DeviceType { return b.kind }
func (b *BaseDevice) State() map[string]any {
	b.mu.RLock()
	defer b.mu.RUnlock()
	// return a shallow copy
	out := make(map[string]any, len(b.state))
	maps.Copy(out, b.state)
	return out
}

func (b *BaseDevice) SetState(k string, v any) {
	b.mu.Lock()
	b.state[k] = v
	b.mu.Unlock()
	ev := Event{DeviceAddr: b.addr, Type: "state_change", Payload: map[string]any{k: v}, Time: time.Now()}
	b.broadcast(ev)
}

func (b *BaseDevice) Subscribe(ch chan<- Event) {
	b.subsMu.Lock()
	b.subs[ch] = struct{}{}
	b.subsMu.Unlock()
}
func (b *BaseDevice) Unsubscribe(ch chan<- Event) {
	b.subsMu.Lock()
	delete(b.subs, ch)
	b.subsMu.Unlock()
}
func (b *BaseDevice) broadcast(e Event) {
	b.subsMu.RLock()
	defer b.subsMu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default: // non-blocking delivery; adapter should consume fast or buffer
		}
	}
}

func (b *BaseDevice) Start(ctx context.Context) error {
	if b.running {
		return fmt.Errorf("device %s already running", b.name)
	}
	b.running = true

	scanner := b.parser.Scanner()
	scanner.AddAddress(b.addr)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case adv, ok := <-scanner.Out(b.addr):
				if !ok {
					continue
				}
				values, ok := b.parser.Parse(adv)
				if !ok {
					continue
				}

				for k, v := range values {
					b.SetState(k, v)
				}
			}
		}

	}()
	return nil
}

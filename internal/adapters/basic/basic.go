package basic

import (
	"context"
	"fmt"
	"gohome/internal/core"
)

type BasicAdapter struct {
	*core.BaseAdapter
	started bool
}

func NewBasicAdapter() core.Adapter {
	h := &BasicAdapter{}
	base := core.NewBaseAdapter("Basic Adapter", h)
	h.BaseAdapter = base
	return h
}

func (h *BasicAdapter) OnStart(ctx context.Context) error {
	h.started = true
	return nil
}

func (h *BasicAdapter) OnStop() error {
	h.started = false
	return nil
}

func (h *BasicAdapter) OnRegisterDevice(d core.Device) (any, context.CancelFunc, error) {
	ctx := h.Context()
	var cancel context.CancelFunc

	switch d.Type() {
	case core.TemperatureSensorType:
		deviceCtx, c := context.WithCancel(ctx)
		cancel = c

		go func() {
			events := make(chan core.Event, 8)
			d.Subscribe(events)
			defer d.Unsubscribe(events)

			for {
				select {
				case ev := <-events:
					if ev.Type == "state_change" {
						if v, ok := ev.Payload["temperature"].(float64); ok {
							if h.started {
								fmt.Printf("[Adapater] (%d:%d:%d) Temperature: %f\n", ev.Time.Hour(), ev.Time.Minute(), ev.Time.Second(), v)
							}
						}
					}
				case <-deviceCtx.Done():
					return
				}
			}
		}()
	case core.HumiditySensorType:
		deviceCtx, c := context.WithCancel(ctx)
		cancel = c

		go func() {
			events := make(chan core.Event, 8)
			d.Subscribe(events)
			defer d.Unsubscribe(events)

			for {
				select {
				case ev := <-events:
					if ev.Type == "state_change" {
						if v, ok := ev.Payload["humidity"].(float64); ok {
							if h.started {
								fmt.Printf("[Adapater] (%d:%d:%d) Humidity: %f\n", ev.Time.Hour(), ev.Time.Minute(), ev.Time.Second(), v)
							}
						}
					}
				case <-deviceCtx.Done():
					return
				}
			}
		}()
	default:
		return nil, nil, fmt.Errorf("unsupported device type: %s", d.Type())
	}

	return nil, cancel, nil
}

func (h *BasicAdapter) OnUnregisterDevice(meta any) error {
	return nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gohome/shared/events"
	"gohome/shared/middlewares"
	"gohome/shared/types"
	"log"
	"net/http"
	"sync"
	"time"
)

type TimeValue struct {
	Timestamp time.Time `json:"x"`
	Value     any       `json:"y"`
}

type HistoryAdapter struct {
	server        *http.Server
	mu            sync.Mutex
	devices       map[string]map[types.CapabilityType][]*TimeValue
	onStateChange func(state types.State)
}

func NewHistoryAdapter(eventBus *events.EventBus, onStateChange func(state types.State), port int) (*HistoryAdapter, error) {
	a := &HistoryAdapter{
		devices:       make(map[string]map[types.CapabilityType][]*TimeValue),
		onStateChange: onStateChange,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/history/device/{deviceId}/capabilities/{capabilityType}", a.getDeviceCapabilityHistory)

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middlewares.CorsMiddleware(mux),
	}

	if err := events.Subscribe(eventBus, events.UpdateDataForAdapter(p.ID), a.onDeviceData); err != nil {
		return nil, err
	}

	if err := events.Subscribe(eventBus, events.RegisterDeviceForAdapter(p.ID), a.onDeviceRegistered); err != nil {
		return nil, err
	}

	if err := events.Subscribe(eventBus, events.UnregisterDeviceForAdapter(p.ID), a.onDeviceUnregistered); err != nil {
		return nil, err
	}

	return a, nil
}

func (h *HistoryAdapter) Start() error {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while running server : %v", err)
			return
		}
	}()

	h.onStateChange(types.StateRunning)

	return nil
}

func (h *HistoryAdapter) Stop() error {
	if err := h.server.Shutdown(context.Background()); err != nil {
		return err
	}
	h.onStateChange(types.StateStopped)
	return nil
}

func (a *HistoryAdapter) onDeviceData(data types.DeviceStateUpdate) {
	a.mu.Lock()
	defer a.mu.Unlock()

	capabilities, ok := a.devices[data.DeviceID]
	if !ok {
		log.Printf("Device with id : %s not registered", data.DeviceID)
		return
	}

	_, ok = capabilities[data.CapabilityType]

	if !ok {
		capabilities[data.CapabilityType] = make([]*TimeValue, 0)
	}

	values := capabilities[data.CapabilityType]

	if len(values) == 0 {
		values = append(values, &TimeValue{
			Value:     data.Value,
			Timestamp: data.Timestamp,
		})
		capabilities[data.CapabilityType] = values
		return
	}

	last := values[len(values)-1]

	if last.Value != data.Value {
		values = append(values, &TimeValue{
			Value:     data.Value,
			Timestamp: data.Timestamp,
		})
	}
	capabilities[data.CapabilityType] = values
}

func (a *HistoryAdapter) onDeviceRegistered(dev types.Device) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, ok := a.devices[dev.ID]
	if ok {
		log.Printf("Device with id : %s already registered", dev.ID)
	}

	a.devices[dev.ID] = make(map[types.CapabilityType][]*TimeValue)
}

func (a *HistoryAdapter) onDeviceUnregistered(dev types.Device) {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, ok := a.devices[dev.ID]
	if !ok {
		log.Printf("Device with id : %s not registered", dev.ID)
	}

	delete(a.devices, dev.ID)
}

func (a *HistoryAdapter) getDeviceCapabilityHistory(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("deviceId")
	capabilityType := types.CapabilityType(r.PathValue("capabilityType"))

	a.mu.Lock()
	values, ok := a.devices[deviceId][capabilityType]
	a.mu.Unlock()

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"x_label": "timestamp",
			"y_label": "values",
			"points":  []*TimeValue{},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"x_label": "timestamp",
		"y_label": "values",
		"points":  values,
	})
}

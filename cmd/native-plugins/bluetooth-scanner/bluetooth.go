package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

const (
	EVENT_DEVICE_FOUND_INTERVAL = time.Minute
	EVENT_DATA_INTERVAL         = time.Minute
)

type BluetoothScanner struct {
	eventBus      *events.EventBus
	adapter       *bluetooth.Adapter
	onStateChange func(state types.State)
	started       bool

	cache map[string]*deviceCacheEntry
	mu    sync.RWMutex
}

type deviceCacheEntry struct {
	lastData     []byte
	lastSeen     time.Time
	lastMetaSent time.Time
}

func NewBluetoothScanner(eventBus *events.EventBus, onStateChange func(state types.State)) *BluetoothScanner {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		log.Fatalf("failed to enable bluetooth adapter: %v", err)
	}
	return &BluetoothScanner{
		eventBus:      eventBus,
		adapter:       adapter,
		onStateChange: onStateChange,
		started:       false,
		cache:         make(map[string]*deviceCacheEntry),
	}
}

func (s *BluetoothScanner) Start() error {
	if s.started {
		return fmt.Errorf("bluetooth scanner already running")
	}

	go s.scanLoop()

	s.onStateChange(types.StateRunning)
	s.started = true
	return nil
}

func (s *BluetoothScanner) scanLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Bluetooth Scanner] Panic: %v", r)
			_ = s.adapter.StopScan()
			s.onStateChange(types.StateStopped)
			s.started = false
		}
	}()

	_ = s.adapter.StopScan()
	time.Sleep(500 * time.Millisecond)

	log.Println("[Bluetooth Scanner] Started")

	err := s.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if len(device.ServiceData()) == 0 {
			return
		}

		serviceData := make(types.BluetoothAdvertisement, len(device.ServiceData()))
		for _, svc := range device.ServiceData() {
			if len(svc.Data) == 0 {
				continue
			}
			serviceData[svc.UUID.String()] = svc.Data
		}

		newDataBytes, err := json.Marshal(serviceData)
		if err != nil {
			log.Printf("[Bluetooth Scanner] Failed to marshal: %v", err)
			return
		}

		address := device.Address.String()
		now := time.Now()

		s.mu.Lock()
		entry, exists := s.cache[address]
		if !exists {
			entry = &deviceCacheEntry{}
			s.cache[address] = entry
		}

		shouldSendRawData := false
		shouldSendMeta := false

		if !bytes.Equal(entry.lastData, newDataBytes) || now.Sub(entry.lastSeen) > EVENT_DATA_INTERVAL {
			shouldSendRawData = true
			entry.lastData = newDataBytes
			entry.lastSeen = now
		}

		if now.Sub(entry.lastMetaSent) > EVENT_DEVICE_FOUND_INTERVAL {
			shouldSendMeta = true
			entry.lastMetaSent = now
		}
		s.mu.Unlock()

		if s.eventBus != nil {
			if shouldSendRawData {
				s.eventBus.Publish(events.Event{
					Type: events.RawDataReceived,
					Payload: &types.RawData{
						Address:     address,
						Data:        newDataBytes,
						Timestamp:   now,
						AddressType: types.BLEAddress,
					},
				})
			}

			if shouldSendMeta {
				s.eventBus.Publish(events.Event{
					Type: events.BluetoothDeviceFound,
					Payload: &BluetoothDevice{
						Name:    device.LocalName(),
						Address: address,
					},
				})
			}
		}
	})

	if err != nil {
		log.Printf("[Bluetooth Scanner] Scan error: %v", err)
		s.onStateChange(types.StateStopped)
		s.started = false
	}
}

func (s *BluetoothScanner) Stop() error {
	if !s.started {
		return nil
	}

	err := s.adapter.StopScan()
	if err != nil {
		log.Printf("Error stopping bluetooth scan: %v", err)
	}

	s.started = false
	s.onStateChange(types.StateStopped)

	log.Println("[Bluetooth Scanner] Stopped")
	return err
}

type BluetoothDevice struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

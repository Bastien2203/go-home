package main

import (
	"encoding/json"
	"fmt"
	"gohome/shared/events"
	"gohome/shared/types"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

type BluetoothScanner struct {
	eventBus      *events.EventBus
	adapter       *bluetooth.Adapter
	onStateChange func(state types.State)
	started       bool
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
			s.onStateChange(types.StateStopped)
			s.started = false
		}
	}()

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

		// log.Printf("[Bluetooth Scanner] Device found: %s (%s) - RSSI: %d dBm",
		// 	device.Address.String(),
		// 	device.LocalName(),
		// 	device.RSSI)

		if s.eventBus != nil {
			bytes, err := json.Marshal(serviceData)
			if err != nil {
				log.Printf("[Bluetooth Scanner] Failed to marshal: %v", err)
			}
			s.eventBus.Publish(events.Event{
				Type: events.RawDataReceived,
				Payload: &types.RawData{
					Address:     device.Address.String(),
					Data:        bytes,
					Timestamp:   time.Now(),
					AddressType: types.BLEAddress,
				},
			})

			s.eventBus.Publish(events.Event{
				Type: events.BluetoothDeviceFound,
				Payload: &BluetoothDevice{
					Name:    device.LocalName(),
					Address: device.Address.String(),
				},
			})
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

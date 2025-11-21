package scanners

import (
	"context"
	"fmt"
	"gohome/internal/core"
	"gohome/internal/events"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

type BluetoothScanner struct {
	eventBus *events.EventBus
	adapter  *bluetooth.Adapter
	started  bool
}

func NewBluetoothScanner(eventBus *events.EventBus) core.Scanner {
	return &BluetoothScanner{
		eventBus: eventBus,
		adapter:  bluetooth.DefaultAdapter,
		started:  false,
	}
}

func (s *BluetoothScanner) ID() string {
	return "bluetooth_scanner"
}

func (s *BluetoothScanner) Name() string {
	return "Bluetooth Scanner"
}

func (s *BluetoothScanner) State() core.State {
	switch s.started {
	case true:
		return core.StateRunning
	default:
		return core.StateStopped
	}
}

func (s *BluetoothScanner) Start(ctx context.Context) error {
	if s.started {
		return fmt.Errorf("bluetooth scanner already running")
	}

	err := s.adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable bluetooth adapter: %v", err)
	}

	go s.scanLoop(ctx)

	s.started = true
	return nil

}

func (s *BluetoothScanner) scanLoop(_ context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Bluetooth Scanner] Panic: %v", r)
			s.started = false
		}
	}()

	log.Println("[Bluetooth Scanner] Started")

	err := s.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if len(device.ServiceData()) == 0 {
			return
		}

		serviceData := make(core.BluetoothAdvertisement, len(device.ServiceData()))
		for _, svc := range device.ServiceData() {
			if len(svc.Data) == 0 {
				continue
			}
			serviceData[svc.UUID] = svc.Data
		}

		// log.Printf("[Bluetooth Scanner] Device found: %s (%s) - RSSI: %d dBm",
		// 	device.Address.String(),
		// 	device.LocalName(),
		// 	device.RSSI)

		if s.eventBus != nil {
			s.eventBus.Publish(events.Event{
				Type: events.RawDataReceived,
				Payload: &core.RawData{
					Address:     device.Address.String(),
					Data:        serviceData,
					Timestamp:   time.Now(),
					AddressType: core.BLEAddress,
				},
			})
		}
	})

	if err != nil {
		log.Printf("[Bluetooth Scanner] Scan error: %v", err)
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

	log.Println("[Bluetooth Scanner] Stopped")
	return err
}

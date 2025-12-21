package main

import (
	"bluetooth-scanner/utils"
	"fmt"
	"log"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

const (
	EVENT_DEVICE_FOUND_INTERVAL = time.Minute
	EVENT_DATA_INTERVAL         = time.Minute
)

type deviceCacheEntry struct {
	address  string
	packetId uint8
}

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

		data := make([]*types.Capability, 0, len(device.ServiceData()))
		protocols := make([]string, 0)
		for _, svc := range device.ServiceData() {
			protocol, ok := ProtocolList[svc.UUID]
			if !ok {
				continue
			}

			protocols = append(protocols, protocol.Name())

			if len(svc.Data) == 0 {
				continue
			}

			if !protocol.CanParse() {
				continue
			}

			capabilities, err := protocol.Parse(device.Address.String(), svc.Data)
			if err != nil {
				if err == utils.DeduplicateBluetoothPackets {
					return
				}
				log.Printf("Error parsing service data for UUID %s: %v", svc.UUID.String(), err)
				continue
			}

			data = append(data, capabilities...)
		}

		for _, mData := range device.ManufacturerData() {
			ManufacturerProtocols, ok := ManufacturerProtocols[mData.CompanyID]
			if !ok {
				continue
			}

			protocols = append(protocols, ManufacturerProtocols.Name())

			if len(mData.Data) == 0 {
				continue
			}

			if !ManufacturerProtocols.CanParse() {
				continue
			}
		}

		address := device.Address.String()
		now := time.Now()

		if s.eventBus != nil {
			if len(data) > 0 {
				s.eventBus.Publish(events.Event{
					Type: events.ParsedDataReceived,
					Payload: &types.ParsedData{
						Address:     address,
						Data:        data,
						Timestamp:   now,
						AddressType: types.BLEAddress,
					},
				})
			}

			s.eventBus.Publish(events.Event{
				Type: events.BluetoothDeviceFound,
				Payload: &BluetoothDevice{
					Name:      device.LocalName(),
					Address:   address,
					Protocols: protocols,
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
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Protocols []string `json:"protocols"`
}

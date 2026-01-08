package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

type BluetoothScanner struct {
	eventBus      *events.EventBus
	adapter       *bluetooth.Adapter
	onStateChange func(state types.State)
	started       bool
	scanResults   chan bluetooth.ScanResult
	lastSeen      map[string]time.Time
	mu            sync.Mutex
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

	s.scanResults = make(chan bluetooth.ScanResult, 100)
	go s.processResults()
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

	err := s.adapter.Scan(func(adapter *bluetooth.Adapter, r bluetooth.ScanResult) {
		select {
		case s.scanResults <- r:
		default:
		}
	})

	if err != nil {
		log.Printf("[Bluetooth Scanner] Scan error: %v", err)
		s.onStateChange(types.StateStopped)
		s.started = false
	}
}

func (s *BluetoothScanner) processResults() {
	lastSeenDevices := make(map[string]time.Time, 100)
	timestamp := time.Now()
	ttl := 1 * time.Hour

	for device := range s.scanResults {
		pData := types.ParsedData{
			Address:     device.Address.String(),
			Timestamp:   time.Now(),
			Data:        make([]*types.Capability, 0, 10),
			AddressType: types.BLEAddress,
		}

		protocolsSeen := make([]string, 0, 5)

		if pData.Timestamp.Sub(timestamp) > ttl {
			lastSeenDevices = make(map[string]time.Time, 100)
			timestamp = pData.Timestamp
		}

		for _, svc := range device.ServiceData() {
			if protocol, ok := ProtocolList[svc.UUID]; ok {
				capabilities := processPayload(svc.Data, protocol, pData.Address)
				protocolsSeen = append(protocolsSeen, protocol.Name())
				pData.Data = append(pData.Data, capabilities...)
			}
		}

		for _, mData := range device.ManufacturerData() {
			if protocol, ok := ManufacturerProtocols[mData.CompanyID]; ok {
				capabilities := processPayload(mData.Data, protocol, pData.Address)
				protocolsSeen = append(protocolsSeen, protocol.Name())
				pData.Data = append(pData.Data, capabilities...)
			}
		}

		if s.eventBus != nil {
			if len(pData.Data) > 0 {
				s.eventBus.Publish(events.Event{
					Type:    events.ParsedDataReceived,
					Payload: pData,
				})
				lastSeen, ok := lastSeenDevices[pData.Address]

				if !ok || pData.Timestamp.Sub(lastSeen) > 30*time.Second {
					lastSeenDevices[pData.Address] = pData.Timestamp
					s.eventBus.Publish(events.Event{
						Type: events.BluetoothDeviceFound,
						Payload: BluetoothDevice{
							Name:      device.LocalName(),
							Address:   pData.Address,
							Protocols: protocolsSeen,
						},
					})
				}
			}
		}
	}
}

func processPayload(payload []byte, protocol Protocol, address string) []*types.Capability {
	if len(payload) == 0 || !protocol.CanParse() {
		return nil
	}

	capabilities, deduplication, err := protocol.Parse(address, payload)
	if err != nil {
		return nil
	}
	if deduplication {
		return nil
	}
	return capabilities
}

func (s *BluetoothScanner) Stop() error {
	if !s.started {
		return nil
	}

	err := s.adapter.StopScan()
	if err != nil {
		log.Printf("Error stopping bluetooth scan: %v", err)
	}

	close(s.scanResults)

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

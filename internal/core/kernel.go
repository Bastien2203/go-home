package core

import (
	"context"
	"fmt"
	"gohome/internal/events"
	"log"
	"sort"
	"sync"
)

type Kernel struct {
	eventBus   *events.EventBus
	repository DeviceRepository
	adapters   map[string]Adapter
	protocols  map[string]Protocol
	scanners   map[string]Scanner
	mu         map[string]*sync.Mutex
	muLock     sync.Mutex
}

func NewKernel(eventBus *events.EventBus, repository DeviceRepository) *Kernel {
	kernel := &Kernel{
		eventBus:   eventBus,
		repository: repository,
		adapters:   make(map[string]Adapter),
		protocols:  make(map[string]Protocol),
		scanners:   make(map[string]Scanner),
		mu:         make(map[string]*sync.Mutex),
	}

	eventBus.Subscribe(events.RawDataReceived, kernel.handleStateUpdate)

	return kernel
}

func (k *Kernel) handleStateUpdate(event events.Event) {
	rawData, ok := event.Payload.(*RawData)
	if !ok {
		log.Printf("[Kernel] (handleStateUpdate) Invalid state update payload")
		return
	}

	device, err := k.repository.FindByAddress(rawData.Address, rawData.AddressType)
	if err != nil || device == nil {
		// log.Printf("[Kernel] (handleStateUpdate) Unknown device for state update: %s", rawData.DeviceID)
		return
	}

	device.LastUpdated = rawData.Timestamp

	protocol, exists := k.protocols[device.Protocol]
	if !exists {
		log.Printf("[Kernel] (handleStateUpdate) Protocol not found for device %s: %s", device.ID, device.Protocol)
		return
	}

	deviceData, err := protocol.Parse(rawData.Data)
	if err != nil {
		log.Printf("[Kernel] (handleStateUpdate) Error parsing data for device %s: %v", device.ID, err)
		return
	}

	mu := k.getMutex(device.ID)
	mu.Lock()
	for _, c := range deviceData {
		device.Capabilities[c.Name] = c
	}

	mu.Unlock()
	if err := k.repository.Save(device); err != nil {
		log.Printf("[Kernel] Error when saving device %s", err.Error())
	}

	for _, adapterID := range device.AdapterIDs {
		if adapter, exists := k.adapters[adapterID]; exists {
			go func(adapter Adapter) {
				for _, c := range deviceData {
					if err := adapter.OnDeviceData(&DeviceStateUpdate{
						DeviceID:       device.ID,
						CapabilityType: c.Name,
						Timestamp:      rawData.Timestamp,
						Value:          c.Value,
					}); err != nil {
						log.Printf("[Kernel] (handleStateUpdate) Error sending data to adapter %s: %v", adapter.ID(), err)
					}
				}
			}(adapter)
		}
	}
}

func (k *Kernel) getMutex(deviceID string) *sync.Mutex {
	k.muLock.Lock()
	defer k.muLock.Unlock()
	if k.mu[deviceID] == nil {
		k.mu[deviceID] = &sync.Mutex{}
	}
	return k.mu[deviceID]
}

func (k *Kernel) Stop() {
	k.StopScanners()
	k.StopAdapters()
}

// --- Scanners Management ---

func (k *Kernel) RegisterScanner(scanner Scanner, ctx context.Context) error {
	k.scanners[scanner.ID()] = scanner

	if err := scanner.Start(ctx); err != nil {
		return err
	}
	log.Printf("[Kernel] Scanner registered: %s", scanner.ID())
	return nil
}

func (k *Kernel) ListScanners() []map[string]any {
	list := make([]map[string]any, 0, len(k.scanners))
	for _, scanner := range k.scanners {
		list = append(list, map[string]any{
			"id":    scanner.ID(),
			"name":  scanner.Name(),
			"state": scanner.State(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["id"].(string) < list[j]["id"].(string)
	})

	return list
}

func (k *Kernel) StopScanners() {
	for _, scanner := range k.scanners {
		if err := scanner.Stop(); err != nil {
			log.Printf("[Kernel] Error stopping scanner %s: %v", scanner.ID(), err)
		} else {
			log.Printf("[Kernel] Scanner stopped: %s", scanner.ID())
		}
	}
}

// --- Adapters Management ---

func (ds *Kernel) RegisterAdapter(adapter Adapter) error {
	if err := adapter.Start(); err != nil {
		return err
	}
	ds.adapters[adapter.ID()] = adapter
	log.Printf("[Kernel] Adapter registered: %s", adapter.ID())
	return nil
}

func (k *Kernel) ListAdapters() []map[string]any {
	list := make([]map[string]any, 0, len(k.adapters))
	for _, adapter := range k.adapters {
		list = append(list, map[string]any{
			"id":    adapter.ID(),
			"name":  adapter.Name(),
			"state": adapter.State(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["id"].(string) < list[j]["id"].(string)
	})
	return list
}

func (k *Kernel) StopAdapters() {
	for _, adapter := range k.adapters {
		if err := adapter.Stop(); err != nil {
			log.Printf("[Kernel] Error stopping adapter %s: %v", adapter.ID(), err)
		} else {
			log.Printf("[Kernel] Adapter stopped: %s", adapter.ID())
		}
	}
}

// --- Protocols Management ---

func (k *Kernel) RegisterProtocol(protocol Protocol) {
	k.protocols[protocol.ID()] = protocol
	log.Printf("[Kernel] Protocol registered: %s", protocol.ID())
}

func (k *Kernel) ListProtocols() []map[string]any {
	list := make([]map[string]any, 0, len(k.protocols))
	for _, p := range k.protocols {
		list = append(list, map[string]any{
			"id":   p.ID(),
			"name": p.Name(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["id"].(string) < list[j]["id"].(string)
	})
	return list
}

func (k *Kernel) GetProtocol(id string) (Protocol, error) {
	protocol, exists := k.protocols[id]
	if !exists {
		return nil, fmt.Errorf("protocol not found: %s", id)
	}
	return protocol, nil
}

// --- Devices Management ---

func (k *Kernel) RegisterDevice(device *Device) error {
	// Sauvegarde initiale
	if err := k.repository.Save(device); err != nil {
		return err
	}
	k.getMutex(device.ID)

	log.Printf("[Kernel] Device registered: %s (ID: %s)", device.Name, device.ID)

	// Auto-link
	for _, adapterID := range device.AdapterIDs {
		if err := k.LinkDeviceToAdapter(device.ID, adapterID); err != nil {
			log.Printf("[Kernel] Warning: Failed to link adapter %s: %v", adapterID, err)
		}
	}
	return nil
}

func (ds *Kernel) GetDevice(deviceID string) (*Device, error) {
	return ds.repository.FindByID(deviceID)
}

func (ds *Kernel) ListDevices() ([]*Device, error) {
	return ds.repository.FindAll()
}

// --- Linking Logic ---

func (k *Kernel) LinkDeviceToAdapter(deviceID, adapterID string) error {
	device, err := k.repository.FindByID(deviceID)
	if err != nil || device == nil {
		return fmt.Errorf("device not found: %s", deviceID)
	}
	adapter, exists := k.adapters[adapterID]
	if !exists {
		return fmt.Errorf("adapter not found: %s", adapterID)
	}

	if err := k.repository.LinkAdapter(deviceID, adapterID); err != nil {
		return err
	}

	return adapter.OnDeviceRegistered(device)
}

func (k *Kernel) UnlinkDeviceFromAdapter(deviceID, adapterID string) error {
	device, err := k.repository.FindByID(deviceID)
	if err != nil || device == nil {
		return fmt.Errorf("device not found: %s", deviceID)
	}
	adapter, exists := k.adapters[adapterID]
	if !exists {
		return fmt.Errorf("adapter not found: %s", adapterID)
	}

	if err := k.repository.UnlinkAdapter(deviceID, adapterID); err != nil {
		return err
	}

	return adapter.OnDeviceUnregistered(device)
}

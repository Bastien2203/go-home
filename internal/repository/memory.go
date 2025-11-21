package repository

import (
	"gohome/internal/core"
	"slices"
	"sync"
)

type InMemoryDeviceRepository struct {
	devices map[string]*core.Device
	mu      sync.RWMutex
}

func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*core.Device),
	}
}

func (r *InMemoryDeviceRepository) Save(device *core.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.devices[device.ID] = device
	return nil
}

func (r *InMemoryDeviceRepository) FindByID(id string) (*core.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, nil
	}
	return device, nil
}

func (r *InMemoryDeviceRepository) FindAll() ([]*core.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*core.Device, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *InMemoryDeviceRepository) LinkAdapter(deviceID, adapterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[deviceID]
	if !exists {
		return nil
	}

	// Check if already linked
	if slices.Contains(device.AdapterIDs, adapterID) {
		return nil
	}

	device.AdapterIDs = append(device.AdapterIDs, adapterID)
	return nil
}

func (r *InMemoryDeviceRepository) UnlinkAdapter(deviceID, adapterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[deviceID]
	if !exists {
		return nil
	}

	device.AdapterIDs = slices.DeleteFunc(device.AdapterIDs, func(e string) bool {
		return e == adapterID
	})

	return nil
}

func (r *InMemoryDeviceRepository) FindByAddress(address string, addressType core.AddressType) (*core.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, device := range r.devices {
		if device.Address == address && device.AddressType == addressType {
			return device, nil
		}
	}
	return nil, nil
}

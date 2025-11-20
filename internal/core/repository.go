package core

import (
	"slices"
	"sync"
)

type InMemoryDeviceRepository struct {
	devices map[string]*Device
	mu      sync.RWMutex
}

func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*Device),
	}
}

func (r *InMemoryDeviceRepository) Save(device *Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.devices[device.ID] = device
	return nil
}

func (r *InMemoryDeviceRepository) FindByID(id string) (*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, nil
	}
	return device, nil
}

func (r *InMemoryDeviceRepository) FindAll() ([]*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*Device, 0, len(r.devices))
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

func (r *InMemoryDeviceRepository) FindByAddress(address string, addressType AddressType) (*Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, device := range r.devices {
		if device.Address == address && device.AddressType == addressType {
			return device, nil
		}
	}
	return nil, nil
}

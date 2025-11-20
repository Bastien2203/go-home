package core

import (
	"time"

	"github.com/google/uuid"
)

type DeviceType string

const (
	TemperatureSensor DeviceType = "thermometer"
)

type RawData struct {
	Address     string      `json:"address"`
	AddressType AddressType `json:"address_type"`
	Data        any         `json:"data"`
	Timestamp   time.Time   `json:"timestamp"`
}

type DeviceStateUpdate struct {
	DeviceID       string         `json:"device_id"`
	CapabilityType CapabilityType `json:"capability_type"`
	Timestamp      time.Time      `json:"timestamp"`
	Value          any            `json:"value"`
}

type Device struct {
	ID           string                         `json:"id"`
	Address      string                         `json:"address"`
	AddressType  AddressType                    `json:"address_type"`
	Name         string                         `json:"name"`
	Protocol     string                         `json:"protocol"`
	AdapterIDs   []string                       `json:"adapter_ids"`
	CreatedAt    time.Time                      `json:"created_at"`
	Capabilities map[CapabilityType]*Capability `json:"capabilities"`
}

func NewDevice(address, name, protocol string, adapterIDs []string, addressType AddressType) *Device {
	return &Device{
		ID:           uuid.New().String(),
		Address:      address,
		AddressType:  addressType,
		Name:         name,
		Protocol:     protocol,
		AdapterIDs:   adapterIDs,
		CreatedAt:    time.Now(),
		Capabilities: make(map[CapabilityType]*Capability),
	}
}

type AddressType string

const (
	BLEAddress AddressType = "ble"
)

type DeviceRepository interface {
	Save(device *Device) error
	FindByID(id string) (*Device, error)
	FindAll() ([]*Device, error)
	LinkAdapter(deviceID, adapterID string) error
	FindByAddress(address string, addressType AddressType) (*Device, error)
}

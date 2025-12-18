package types

import (
	"time"

	"github.com/google/uuid"
)

type DeviceType string

const (
	TemperatureSensor DeviceType = "thermometer"
)

type ParsedData struct {
	Address     string        `json:"address"`
	AddressType AddressType   `json:"address_type"`
	Data        []*Capability `json:"data"`
	Timestamp   time.Time     `json:"timestamp"`
}

type DeviceStateUpdate struct {
	DeviceID       string         `json:"device_id"`
	DeviceName     string         `json:"name"`
	CapabilityType CapabilityType `json:"capability_type"`
	Timestamp      time.Time      `json:"timestamp"`
	Value          any            `json:"value"`
}

type Device struct {
	ID           string                         `json:"id"`
	Address      string                         `json:"address"`
	AddressType  AddressType                    `json:"address_type"`
	Name         string                         `json:"name"`
	AdapterIDs   []string                       `json:"adapter_ids"`
	CreatedAt    time.Time                      `json:"created_at"`
	Capabilities map[CapabilityType]*Capability `json:"capabilities"`
	LastUpdated  time.Time                      `json:"last_updated"`
}

func NewDevice(address, name string, adapterIDs []string, addressType AddressType) *Device {
	return &Device{
		ID:           uuid.New().String(),
		Address:      address,
		AddressType:  addressType,
		Name:         name,
		AdapterIDs:   adapterIDs,
		CreatedAt:    time.Now(),
		Capabilities: make(map[CapabilityType]*Capability),
		LastUpdated:  time.Now(),
	}
}

type AddressType string

const (
	BLEAddress   AddressType = "ble"
	BasicAddress AddressType = "basic"
)

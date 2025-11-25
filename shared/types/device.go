package types

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
	Data        []byte      `json:"data"`
	Timestamp   time.Time   `json:"timestamp"`
}

type DeviceStateUpdate struct {
	DeviceID       string         `json:"device_id"`
	DeviceName     string         `json:"name"`
	DeviceProtocol string         `json:"protocol"`
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
	LastUpdated  time.Time                      `json:"last_updated"`
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
		LastUpdated:  time.Now(),
	}
}

type AddressType string

const (
	BLEAddress   AddressType = "ble"
	BasicAddress AddressType = "basic"
)

package types

import (
	"time"

	"github.com/absmach/senml"
	"github.com/google/uuid"
)

type DeviceType string

const (
	TemperatureSensor DeviceType = "thermometer"
)

type ParsedData struct {
	Address     string      `json:"address" cbor:"1,keyasint"`
	AddressType AddressType `json:"address_type" cbor:"2,keyasint"`
	Data        senml.Pack  `json:"data" cbor:"3,keyasint"`
	Timestamp   time.Time   `json:"timestamp" cbor:"4,keyasint"`
}

type DeviceStateUpdate struct {
	DeviceID   string     `json:"device_id" cbor:"1,keyasint"`
	DeviceName string     `json:"name" cbor:"2,keyasint"`
	Data       senml.Pack `json:"data" cbor:"3,keyasint"`
}

type Device struct {
	ID          string      `json:"id"`
	Address     string      `json:"address"`
	AddressType AddressType `json:"address_type"`
	Name        string      `json:"name"`
	AdapterIDs  []string    `json:"adapter_ids"`
	CreatedAt   time.Time   `json:"created_at"`
	LastUpdated time.Time   `json:"last_updated"`
}

func NewDevice(address, name string, adapterIDs []string, addressType AddressType) *Device {
	return &Device{
		ID:          uuid.New().String(),
		Address:     address,
		AddressType: addressType,
		Name:        name,
		AdapterIDs:  adapterIDs,
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	}
}

type AddressType string

const (
	BLEAddress   AddressType = "ble"
	BasicAddress AddressType = "basic"
)

package core

import "github.com/Bastien2203/go-home/shared/types"

type DeviceRepository interface {
	Save(device *types.Device) error
	FindByID(id string) (*types.Device, error)
	FindAll() ([]*types.Device, error)
	LinkAdapter(deviceID, adapterID string) error
	UnlinkAdapter(deviceID, adapterID string) error
	FindByAddress(address string, addressType types.AddressType) (*types.Device, error)
	Delete(deviceId string) error
}

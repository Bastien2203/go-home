package core

type Adapter interface {
	ID() string
	Name() string
	Start() error
	Stop() error
	OnDeviceData(data *DeviceStateUpdate) error
	OnDeviceRegistered(device *Device) error
	State() State
}

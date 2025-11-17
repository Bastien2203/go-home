package state

import "gohome/internal/core"

type Store interface {
	GetAdapters() map[string]core.Adapter
	GetDevices() map[string]core.Device
	GetParsers() map[string]core.Parser

	GetDevice(addr string) (core.Device, bool)
	GetAdapter(id string) (core.Adapter, bool)

	SaveDevice(device core.Device)
	SaveAdapter(adapter core.Adapter)
	SaveParser(parser core.Parser)

	DeleteDevice(addr string)
}

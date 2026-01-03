package main

import (
	"bluetooth-scanner/protocols"

	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

type Protocol interface {
	Name() string
	Parse(address string, payload []byte) ([]*types.Capability, error)
	CanParse() bool
}

var ProtocolList = map[bluetooth.UUID]Protocol{
	protocols.BthomeUUID:           protocols.NewBthomeParser(),
	bluetooth.New16BitUUID(0x181C): protocols.NewBthomeParser(),
	bluetooth.New16BitUUID(0xFE95): protocols.NewNotImplementedParser("Xiaomi / Mijia"),
	bluetooth.New16BitUUID(0xFD3D): protocols.NewSwitchBotParser(),
	bluetooth.New16BitUUID(0xFEAA): protocols.NewNotImplementedParser("Eddystone (Google)"),
	bluetooth.New16BitUUID(0xFEED): protocols.NewNotImplementedParser("Tile"),
	bluetooth.New16BitUUID(0xFE2C): protocols.NewNotImplementedParser("Google"),
	bluetooth.New16BitUUID(0xFD6F): protocols.NewNotImplementedParser("Exposure Notification"),
}

var ManufacturerProtocols = map[uint16]Protocol{
	0x004C: protocols.NewNotImplementedParser("Apple"),
	0x0059: protocols.NewNotImplementedParser("Nordic Semiconductor"),
	0x0499: protocols.NewNotImplementedParser("Ruuvi Innovations"),
	0xEC88: protocols.NewNotImplementedParser("Govee"),
	0x0075: protocols.NewNotImplementedParser("Samsung"),
	0x0006: protocols.NewNotImplementedParser("Microsoft"),
	0x0157: protocols.NewNotImplementedParser("Anker"),
	0x0969: protocols.NewNotImplementedParser("SwitchBot"),
}

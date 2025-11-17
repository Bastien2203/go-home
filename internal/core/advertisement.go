package core

import (
	"time"

	"tinygo.org/x/bluetooth"
)

type Advertisment interface {
	IsAdvertisment()
}

type BasicAdvertisment struct {
	Value float64
	Type  string
}

func (a *BasicAdvertisment) IsAdvertisment() {}

type BluetoothAdvertisement struct {
	Addr        string // mac
	Rssi        int16
	Time        time.Time
	ServiceData map[bluetooth.UUID][]byte
}

func (a *BluetoothAdvertisement) IsAdvertisment() {}

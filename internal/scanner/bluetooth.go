package scanner

import (
	"context"
	"gohome/internal/core"

	"tinygo.org/x/bluetooth"
)

type BluetoothScanner struct {
	adapter *bluetooth.Adapter
	core.BaseScanner
}

func NewBluetoothScanner(addresses []string) *BluetoothScanner {
	this := &BluetoothScanner{
		adapter: bluetooth.DefaultAdapter,
	}

	this.BaseScanner = *core.NewBaseScanner(addresses, this)
	return this
}

func (s *BluetoothScanner) OnStart(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.adapter.StopScan()
	}()

	err := s.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		var sd = make(map[bluetooth.UUID][]byte, len(device.ServiceData()))
		for _, svc := range device.ServiceData() {
			sd[svc.UUID] = svc.Data
		}

		adv := &core.BluetoothAdvertisement{
			Addr:        device.Address.String(),
			Rssi:        device.RSSI,
			ServiceData: sd,
		}

		select {
		case s.BaseScanner.ChOut[adv.Addr] <- adv:
		default:
			//log.Printf("bluetooth scanner queue full, dropping advertisement from %s\n", device.Address.String())
		}

	})

	return err
}

func (s *BluetoothScanner) OnStop() error {
	return s.adapter.StopScan()
}

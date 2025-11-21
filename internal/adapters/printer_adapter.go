package adapters

import (
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"time"
)

type PrinterAdapter struct {
	kernel       *core.Kernel
	adapterState core.State
}

func NewPrinterAdapter(k *core.Kernel) *PrinterAdapter {
	return &PrinterAdapter{
		kernel:       k,
		adapterState: core.StateStopped,
	}
}

func (p *PrinterAdapter) ID() string {
	return "printer"
}

func (p *PrinterAdapter) State() core.State {
	return p.adapterState
}

func (h *PrinterAdapter) Name() string {
	return "Printer"
}

func (p *PrinterAdapter) Start() error {
	fmt.Println("[Printer adapter] started")
	p.adapterState = core.StateRunning
	return nil
}

func (p *PrinterAdapter) Stop() error {
	fmt.Println("[Printer adapter] stopped")
	p.adapterState = core.StateStopped
	return nil
}

func (p *PrinterAdapter) OnDeviceData(data *core.DeviceStateUpdate) error {
	device, err := p.kernel.GetDevice(data.DeviceID)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Device Data Received:\n%s\n",
		time.Now().Format("15:04:05"),
		string(jsonData))
	return nil
}

func (p *PrinterAdapter) OnDeviceRegistered(dev *core.Device) error {
	fmt.Printf("[Printer Adapter] New device registered: %s (%s)\n", dev.Name, dev.ID)
	return nil
}

func (p *PrinterAdapter) OnDeviceUnregistered(dev *core.Device) error {
	fmt.Printf("[Printer Adapter] device unregistered: %s (%s)\n", dev.Name, dev.ID)
	return nil
}

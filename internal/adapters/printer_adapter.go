package adapters

import (
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"time"
)

type PrinterAdapter struct {
	kernel *core.Kernel
}

func NewPrinterAdapter(k *core.Kernel) *PrinterAdapter {
	return &PrinterAdapter{
		kernel: k,
	}
}

func (p *PrinterAdapter) ID() string {
	return "printer"
}

func (h *PrinterAdapter) Name() string {
	return "Printer"
}

func (p *PrinterAdapter) Start() error {
	fmt.Println("Printer adapter started")
	return nil
}

func (p *PrinterAdapter) Stop() error {
	fmt.Println("Printer adapter stopped")
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
	fmt.Printf("Printer Adapter: New device registered: %s (%s)\n", dev.Name, dev.ID)
	return nil
}

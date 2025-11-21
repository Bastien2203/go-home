package adapters

import (
	"context"
	"fmt"
	"log"
	"time"

	"gohome/internal/core"
	"gohome/utils"
	"sync"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type HomekitAdapter struct {
	kernel *core.Kernel

	// State Management
	accessories map[string]*accessory.A
	knownCaps   map[string]map[core.CapabilityType]bool
	bridge      *accessory.Bridge

	// Server & Concurrency Control
	server       *hap.Server
	serverStop   context.CancelFunc
	mu           sync.Mutex
	restartTimer *time.Timer
	adapterState core.State
}

func NewHomeKitAdapter(k *core.Kernel) *HomekitAdapter {
	info := accessory.Info{
		Name:         "GoHome Hub",
		SerialNumber: "GOHOME-HUB-001",
		Manufacturer: "GoHome",
		Model:        "Bridge v2",
		Firmware:     "1.0.0",
	}
	bridge := accessory.NewBridge(info)

	return &HomekitAdapter{
		kernel:       k,
		accessories:  make(map[string]*accessory.A),
		knownCaps:    make(map[string]map[core.CapabilityType]bool),
		bridge:       bridge,
		adapterState: core.StateStopped,
	}
}

func (h *HomekitAdapter) ID() string {
	return "homekit"
}

func (h *HomekitAdapter) State() core.State {
	return h.adapterState
}

func (h *HomekitAdapter) Name() string {
	return "Homekit"
}

func (h *HomekitAdapter) Start() error {
	log.Println("[HomeKit] Starting adapter...")

	// Pre-load existing devices from Kernel
	devices, _ := h.kernel.ListDevices()
	for _, dev := range devices {
		h.updateDeviceStructure(dev.ID, dev.Name, dev.Protocol, dev.Capabilities)
	}

	return h.reloadServer()
}

func (h *HomekitAdapter) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.serverStop != nil {
		h.serverStop()
		h.serverStop = nil
	}
	log.Println("[HomeKit] Adapter stopped.")
	return nil
}

func (h *HomekitAdapter) OnDeviceData(data *core.DeviceStateUpdate) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, isSupported := capabilityToHAP[data.CapabilityType]; !isSupported {
		return nil
	}

	newCap := map[core.CapabilityType]*core.Capability{
		data.CapabilityType: {},
	}

	name := data.DeviceID
	protocol := "unknown"
	if dev, err := h.kernel.GetDevice(data.DeviceID); err == nil && dev != nil {
		name = dev.Name
		protocol = dev.Protocol
	}

	// Check if this capability forces a structure change
	// If so, rebuild accessory and schedule server reload
	if changed := h.updateDeviceStructure(data.DeviceID, name, protocol, newCap); changed {
		log.Printf("[HomeKit] New capability detected: %s for device %s", data.CapabilityType, data.DeviceID)
		h.scheduleReload()
	}

	acc, exists := h.accessories[data.DeviceID]
	if !exists {
		return fmt.Errorf("device %s not found in accessories map", data.DeviceID)
	}

	h.updateAccessoryValue(acc, data)

	return nil
}

func (h *HomekitAdapter) OnDeviceRegistered(dev *core.Device) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[HomeKit] Device registered : %s", dev.ID)

	if changed := h.updateDeviceStructure(dev.ID, dev.Name, dev.Protocol, dev.Capabilities); changed {
		h.scheduleReload()
	}

	return nil
}

func (h *HomekitAdapter) OnDeviceUnregistered(dev *core.Device) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[HomeKit] Device unregistered : %s", dev.ID)

	delete(h.knownCaps, dev.ID)
	delete(h.accessories, dev.ID)

	h.scheduleReload()

	return nil
}

func (h *HomekitAdapter) updateDeviceStructure(deviceID, name, protocol string, capabilities map[core.CapabilityType]*core.Capability) bool {
	if _, exists := h.knownCaps[deviceID]; !exists {
		h.knownCaps[deviceID] = make(map[core.CapabilityType]bool)
	}

	// Detect changes
	hasChanges := false
	hasAtLeastOneSupportedCap := false

	for capType := range capabilities {
		if _, isSupported := capabilityToHAP[capType]; !isSupported {
			continue
		}
		hasAtLeastOneSupportedCap = true
		if !h.knownCaps[deviceID][capType] {
			h.knownCaps[deviceID][capType] = true
			hasChanges = true
		}
	}

	// Check existing capabilities (in case we are updating an existing device)
	if !hasAtLeastOneSupportedCap {
		for capType := range h.knownCaps[deviceID] {
			if _, isSupported := capabilityToHAP[capType]; isSupported {
				hasAtLeastOneSupportedCap = true
				break
			}
		}
	}

	if !hasAtLeastOneSupportedCap {
		return false
	}

	// If changes detected OR device is missing from accessories, rebuild it
	_, accExists := h.accessories[deviceID]
	if !accExists || hasChanges {
		h.accessories[deviceID] = h.rebuildAccessory(deviceID, name, protocol)
		return true
	}

	return false
}

func (h *HomekitAdapter) rebuildAccessory(deviceID, name, protocol string) *accessory.A {
	info := accessory.Info{
		Name:         name,
		SerialNumber: deviceID,
		Manufacturer: "GoHome",
		Model:        protocol,
		Firmware:     "1.0.0",
	}
	acc := accessory.New(info, accessory.TypeSensor)

	if caps, ok := h.knownCaps[deviceID]; ok {
		for capType := range caps {
			if mapping, supported := capabilityToHAP[capType]; supported {
				switch mapping.serviceType {
				case service.TypeTemperatureSensor:
					acc.AddS(service.NewTemperatureSensor().S)
				case service.TypeHumiditySensor:
					acc.AddS(service.NewHumiditySensor().S)
				case service.TypeBatteryService:
					acc.AddS(service.NewBatteryService().S)
				}
			}
		}
	}
	return acc
}

func (h *HomekitAdapter) updateAccessoryValue(acc *accessory.A, data *core.DeviceStateUpdate) {
	mapping, supported := capabilityToHAP[data.CapabilityType]
	if !supported {
		return
	}

	svc := findService(acc, mapping.serviceType)
	if svc == nil {
		return
	}

	char := svc.C(mapping.characteristicType)
	if char != nil {
		updateCharacteristic(char, data.Value)
	}
}

func (h *HomekitAdapter) scheduleReload() {
	h.adapterState = core.StateRestarting
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		if h.restartTimer != nil {
			h.restartTimer.Stop()
		}

		log.Println("[HomeKit] Structure change detected. Scheduling reload in 2s...")

		h.restartTimer = time.AfterFunc(2*time.Second, func() {
			h.mu.Lock()
			defer h.mu.Unlock()

			log.Println("[HomeKit] Executing scheduled server reload...")
			if err := h.reloadServer(); err != nil {
				log.Printf("[HomeKit] Error reloading server: %v", err)
			}
			h.restartTimer = nil
		})
	}()
}

func (h *HomekitAdapter) reloadServer() error {
	h.adapterState = core.StateStopped
	if h.serverStop != nil {
		h.serverStop()
		h.serverStop = nil
		time.Sleep(100 * time.Millisecond)
	}

	var accs []*accessory.A
	for _, a := range h.accessories {
		accs = append(accs, a)
	}

	if len(accs) == 0 {
		log.Println("[HomeKit] No accessories to publish yet.")
		return nil
	}

	fs := hap.NewFsStore("./homekit_data")

	server, err := hap.NewServer(fs, h.bridge.A, accs...)
	if err != nil {
		return fmt.Errorf("failed to create hap server: %w", err)
	}
	h.server = server

	ctx, cancel := context.WithCancel(context.Background())
	h.serverStop = cancel

	go func() {
		log.Printf("[HomeKit] Server running with %d accessories", len(accs))
		if err := server.ListenAndServe(ctx); err != nil && err != context.Canceled {
			log.Printf("[HomeKit] Server stopped with error: %v", err)
		}
	}()

	h.adapterState = core.StateRunning

	return nil
}

// --- Helpers ---

func findService(acc *accessory.A, serviceType string) *service.S {
	for _, s := range acc.Ss {
		if s.Type == serviceType {
			return s
		}
	}
	return nil
}

func updateCharacteristic(c *characteristic.C, val any) {
	if c == nil {
		return
	}
	switch c.Format {
	case characteristic.FormatFloat:
		if v, ok := utils.ToFloat(val); ok {
			(&characteristic.Float{C: c}).SetValue(v)
		}
	case characteristic.FormatInt32, characteristic.FormatUInt8, characteristic.FormatUInt16:
		if v, ok := utils.ToInt(val); ok {
			(&characteristic.Int{C: c}).SetValue(v)
		}
	case characteristic.FormatBool:
		if v, ok := val.(bool); ok {
			(&characteristic.Bool{C: c}).SetValue(v)
		}
	case characteristic.FormatString:
		if v, ok := val.(string); ok {
			(&characteristic.String{C: c}).SetValue(v)
		}
	}
}

// Definition of supported HomeKit mappings
var capabilityToHAP = map[core.CapabilityType]struct {
	characteristicType string
	serviceType        string
}{
	core.CapabilityTemperature: {
		characteristicType: characteristic.TypeCurrentTemperature,
		serviceType:        service.TypeTemperatureSensor,
	},
	core.CapabilityHumidity: {
		characteristicType: characteristic.TypeCurrentRelativeHumidity,
		serviceType:        service.TypeHumiditySensor,
	},
	core.CapabilityBattery: {
		characteristicType: characteristic.TypeBatteryLevel,
		serviceType:        service.TypeBatteryService,
	},
}

package main

import (
	"context"
	"fmt"

	"log"
	"os"
	"sync"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
	"github.com/Bastien2203/go-home/utils"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type HomekitAdapter struct {

	// State Management
	accessories map[string]*accessory.A
	knownCaps   map[string]map[types.CapabilityType]bool
	bridge      *accessory.Bridge

	// Server & Concurrency Control
	server         *hap.Server
	serverStop     context.CancelFunc
	mu             sync.Mutex
	restartTimer   *time.Timer
	adapterState   types.State
	onStateChange  func(state types.State)
	homekitDataDir string
}

func NewHomeKitAdapter(eventBus *events.EventBus, onStateChange func(state types.State), homekitDataDir string) (*HomekitAdapter, error) {
	info := accessory.Info{
		Name:         "GoHome Hub",
		SerialNumber: "GOHOME-HUB-001",
		Manufacturer: "GoHome",
		Model:        "Bridge v2",
		Firmware:     "1.0.0",
	}
	bridge := accessory.NewBridge(info)

	a := &HomekitAdapter{
		accessories:    make(map[string]*accessory.A),
		knownCaps:      make(map[string]map[types.CapabilityType]bool),
		bridge:         bridge,
		adapterState:   types.StateStopped,
		onStateChange:  onStateChange,
		homekitDataDir: homekitDataDir,
	}

	if err := events.Subscribe(eventBus, events.UpdateDataForAdapter(p.ID), a.onDeviceData); err != nil {
		return nil, err
	}

	if err := events.Subscribe(eventBus, events.RegisterDeviceForAdapter(p.ID), a.onDeviceRegistered); err != nil {
		return nil, err
	}

	if err := events.Subscribe(eventBus, events.UnregisterDeviceForAdapter(p.ID), a.onDeviceUnregistered); err != nil {
		return nil, err
	}

	return a, nil
}

func (h *HomekitAdapter) Start() error {
	log.Println("[HomeKit] Starting adapter...")

	// // Pre-load existing devices from Kernel
	// devices, _ := h.kernel.ListDevices()
	// for _, dev := range devices {
	// 	h.updateDeviceStructure(dev.ID, dev.Name, dev.Protocol, dev.Capabilities)
	// }

	return h.reloadServer()
}

func (h *HomekitAdapter) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.serverStop != nil {
		h.serverStop()
		h.onStateChange(types.StateStopped)
		h.serverStop = nil
	}
	log.Println("[HomeKit] Adapter stopped.")
	return nil
}

func (h *HomekitAdapter) onDeviceData(data types.DeviceStateUpdate) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, isSupported := capabilityToHAP[data.CapabilityType]; !isSupported {
		return
	}

	newCap := map[types.CapabilityType]*types.Capability{
		data.CapabilityType: {},
	}

	// Check if this capability forces a structure change
	// If so, rebuild accessory and schedule server reload
	if changed := h.updateDeviceStructure(data.DeviceID, data.DeviceName, newCap); changed {
		log.Printf("[HomeKit] New capability detected: %s for device %s", data.CapabilityType, data.DeviceID)
		h.scheduleReload()
	}

	acc, exists := h.accessories[data.DeviceID]
	if !exists {
		log.Printf("device %s not found in accessories map", data.DeviceID)
		return
	}

	h.updateAccessoryValue(acc, &data)
}

func (h *HomekitAdapter) onDeviceRegistered(dev types.Device) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[HomeKit] Device registered : %s", dev.ID)

	if changed := h.updateDeviceStructure(dev.ID, dev.Name, dev.Capabilities); changed {
		h.scheduleReload()
	}
}

func (h *HomekitAdapter) onDeviceUnregistered(dev types.Device) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("[HomeKit] Device unregistered : %s", dev.ID)

	delete(h.knownCaps, dev.ID)
	delete(h.accessories, dev.ID)

	h.scheduleReload()
}

func (h *HomekitAdapter) updateDeviceStructure(deviceID, name string, capabilities map[types.CapabilityType]*types.Capability) bool {
	if _, exists := h.knownCaps[deviceID]; !exists {
		h.knownCaps[deviceID] = make(map[types.CapabilityType]bool)
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
		h.accessories[deviceID] = h.rebuildAccessory(deviceID, name)
		return true
	}

	return false
}

func (h *HomekitAdapter) rebuildAccessory(deviceID, name string) *accessory.A {
	info := accessory.Info{
		Name:         name,
		SerialNumber: deviceID,
		Manufacturer: "GoHome",
		Model:        "",
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

func (h *HomekitAdapter) updateAccessoryValue(acc *accessory.A, data *types.DeviceStateUpdate) {
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
	h.adapterState = types.StateRestarting
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
	h.adapterState = types.StateStopped
	h.onStateChange(types.StateStopped)
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

	fs := hap.NewFsStore(h.homekitDataDir)

	server, err := hap.NewServer(fs, h.bridge.A, accs...)
	if err != nil {
		return fmt.Errorf("failed to create hap server: %w", err)
	}
	server.Ifaces = []string{os.Getenv("INTERNET_INTERFACE")}
	h.server = server

	ctx, cancel := context.WithCancel(context.Background())
	h.serverStop = cancel

	go func() {
		log.Printf("[HomeKit] Server running with %d accessories", len(accs))
		if err := server.ListenAndServe(ctx); err != nil && err != context.Canceled {
			log.Printf("[HomeKit] Server stopped with error: %v", err)
		}
	}()

	h.adapterState = types.StateRunning
	h.onStateChange(types.StateRunning)

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
var capabilityToHAP = map[types.CapabilityType]struct {
	characteristicType string
	serviceType        string
}{
	types.CapabilityTemperature: {
		characteristicType: characteristic.TypeCurrentTemperature,
		serviceType:        service.TypeTemperatureSensor,
	},
	types.CapabilityHumidity: {
		characteristicType: characteristic.TypeCurrentRelativeHumidity,
		serviceType:        service.TypeHumiditySensor,
	},
	types.CapabilityBattery: {
		characteristicType: characteristic.TypeBatteryLevel,
		serviceType:        service.TypeBatteryService,
	},
	types.CapabilityButtonEvent: {
		characteristicType: characteristic.TypeProgrammableSwitchEvent,
		serviceType:        service.TypeStatelessProgrammableSwitch,
	},
}

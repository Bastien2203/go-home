package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gohome/shared/events"
	"gohome/shared/plugin"
	"gohome/shared/types"

	"log"
	"sort"
	"sync"
)

type Kernel struct {
	eventBus      *events.EventBus
	repository    DeviceRepository
	protocols     map[string]types.Protocol
	mu            map[string]*sync.Mutex
	muLock        sync.Mutex
	pluginManager *PluginManager
	processes     map[string]*exec.Cmd
	pMu           sync.Mutex
}

func NewKernel(eventBus *events.EventBus, repository DeviceRepository) (*Kernel, error) {
	pluginManager, err := NewPluginManager(eventBus)
	if err != nil {
		return nil, err
	}

	kernel := &Kernel{
		eventBus:      eventBus,
		repository:    repository,
		protocols:     make(map[string]types.Protocol),
		mu:            make(map[string]*sync.Mutex),
		processes:     make(map[string]*exec.Cmd),
		pluginManager: pluginManager,
	}

	if err := events.Subscribe(eventBus, events.RawDataReceived, kernel.handleStateUpdate); err != nil {
		return nil, err
	}

	return kernel, nil
}

func (k *Kernel) handleStateUpdate(rawData types.RawData) {
	device, err := k.repository.FindByAddress(rawData.Address, rawData.AddressType)
	if err != nil || device == nil {
		return
	}

	device.LastUpdated = rawData.Timestamp

	protocol, exists := k.protocols[device.Protocol]
	if !exists {
		log.Printf("[Kernel] (handleStateUpdate) Protocol not found for device %s: %s", device.ID, device.Protocol)
		return
	}

	deviceData, err := protocol.Parse(rawData.Data)
	if err != nil {
		log.Printf("[Kernel] (handleStateUpdate) Error parsing data for device %s: %v", device.ID, err)
		return
	}

	mu := k.getMutex(device.ID)
	mu.Lock()
	for _, c := range deviceData {
		device.Capabilities[c.Name] = c
	}

	mu.Unlock()
	// if err := k.repository.Save(device); err != nil {
	// 	log.Printf("[Kernel] Error when saving device %s", err.Error())
	// }

	for _, adapterID := range device.AdapterIDs {
		go func(adapterID string) {
			for _, c := range deviceData {
				k.eventBus.Publish(events.Event{
					Type: events.UpdateDataForAdapter(adapterID),
					Payload: types.DeviceStateUpdate{
						DeviceID:       device.ID,
						DeviceName:     device.Name,
						DeviceProtocol: device.Protocol,
						CapabilityType: c.Name,
						Timestamp:      rawData.Timestamp,
						Value:          c.Value,
					},
				})
			}
		}(adapterID)
	}
}

func (k *Kernel) getMutex(deviceID string) *sync.Mutex {
	k.muLock.Lock()
	defer k.muLock.Unlock()
	if k.mu[deviceID] == nil {
		k.mu[deviceID] = &sync.Mutex{}
	}
	return k.mu[deviceID]
}

func (k *Kernel) deleteMutex(deviceID string) {
	k.muLock.Lock()
	defer k.muLock.Unlock()
	delete(k.mu, deviceID)
}

func (k *Kernel) LoadPlugins(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err == nil {
			if info.Mode()&0111 == 0 {
				log.Printf("Ignore file %s: non executable", f.Name())
				continue
			}
		}

		fileName := f.Name()
		fullPath := filepath.Join(dir, fileName)

		cmd := exec.Command(fullPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		k.pMu.Lock()
		k.processes[fileName] = cmd
		k.pMu.Unlock()

		go func(name string, c *exec.Cmd) {
			defer func() {
				k.pMu.Lock()
				delete(k.processes, name)
				k.pMu.Unlock()
			}()

			log.Printf("[%s] Starting...", name)
			if err := c.Run(); err != nil {
				log.Printf("[%s] Stopped or Error: %v", name, err)
				return
			}
			log.Printf("[%s] Finished naturally", name)
		}(fileName, cmd)
	}
	return nil
}

func (k *Kernel) UnloadPlugins() {
	k.pMu.Lock()
	defer k.pMu.Unlock()

	log.Println("Unloading all plugins...")

	for name, cmd := range k.processes {
		if cmd.Process != nil {
			log.Printf("Signaling %s to stop...", name)
			err := cmd.Process.Signal(os.Interrupt)

			if err != nil {
				log.Printf("Failed to kill %s: %v", name, err)
			}
		}
	}
}

// --- Scanners Management ---

func (k *Kernel) ListScanners() []*plugin.Plugin {
	scanners := k.pluginManager.GetPluginsByType(plugin.PluginScanner)
	sort.Slice(scanners, func(i, j int) bool {
		return scanners[i].ID < scanners[j].ID
	})
	return scanners
}

func (k *Kernel) StopScanner(id string) error {
	plugin, err := k.pluginManager.GetPluginById(plugin.PluginScanner, id)
	if err != nil {
		return err
	}

	return k.pluginManager.StopPlugin(plugin)
}

func (k *Kernel) StartScanner(id string) error {
	plugin, err := k.pluginManager.GetPluginById(plugin.PluginScanner, id)
	if err != nil {
		return err
	}

	return k.pluginManager.StartPlugin(plugin)
}

func (k *Kernel) StopScanners() {
	scanners := k.pluginManager.GetPluginsByType(plugin.PluginScanner)
	for _, scanner := range scanners {
		if err := k.pluginManager.StopPlugin(scanner); err != nil {
			log.Printf("error stoppig scanner : %v", err)
		}
	}
}

// --- Adapters Management ---

func (k *Kernel) ListAdapters() []*plugin.Plugin {
	adapters := k.pluginManager.GetPluginsByType(plugin.PluginAdapter)
	sort.Slice(adapters, func(i, j int) bool {
		return adapters[i].ID < adapters[j].ID
	})
	return adapters
}

func (k *Kernel) StopAdapter(id string) error {
	plugin, err := k.pluginManager.GetPluginById(plugin.PluginAdapter, id)
	if err != nil {
		return err
	}

	return k.pluginManager.StopPlugin(plugin)
}

func (k *Kernel) StartAdapter(id string) error {
	plugin, err := k.pluginManager.GetPluginById(plugin.PluginAdapter, id)
	if err != nil {
		return err
	}

	return k.pluginManager.StartPlugin(plugin)
}

func (k *Kernel) StopAdapters() {
	adapters := k.pluginManager.GetPluginsByType(plugin.PluginAdapter)
	for _, adapter := range adapters {
		if err := k.pluginManager.StopPlugin(adapter); err != nil {
			log.Printf("error stoppig adapter : %v", err)
		}
	}
}

// --- Protocols Management ---

func (k *Kernel) RegisterProtocol(protocol types.Protocol) {
	k.protocols[protocol.ID()] = protocol
	log.Printf("[Kernel] Protocol registered: %s", protocol.ID())
}

func (k *Kernel) ListProtocols() []map[string]any {
	list := make([]map[string]any, 0, len(k.protocols))
	for _, p := range k.protocols {
		list = append(list, map[string]any{
			"id":   p.ID(),
			"name": p.Name(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i]["id"].(string) < list[j]["id"].(string)
	})
	return list
}

func (k *Kernel) GetProtocol(id string) (types.Protocol, error) {
	protocol, exists := k.protocols[id]
	if !exists {
		return nil, fmt.Errorf("protocol not found: %s", id)
	}
	return protocol, nil
}

// --- Devices Management ---

func (k *Kernel) RegisterDevice(device *types.Device) error {
	if err := k.repository.Save(device); err != nil {
		return err
	}
	k.getMutex(device.ID)

	log.Printf("[Kernel] Device registered: %s (ID: %s)", device.Name, device.ID)

	for _, adapterID := range device.AdapterIDs {
		if err := k.LinkDeviceToAdapter(device.ID, adapterID); err != nil {
			log.Printf("[Kernel] Warning: Failed to link adapter %s: %v", adapterID, err)
		}
	}
	return nil
}

func (k *Kernel) UnregisterDevice(deviceId string) error {
	device, err := k.repository.FindByID(deviceId)
	if err != nil {
		return fmt.Errorf("device doesnt exists : %v", err)
	}

	for _, adapterID := range device.AdapterIDs {
		if err := k.UnlinkDeviceFromAdapter(device.ID, adapterID); err != nil {
			log.Printf("[Kernel] Warning: Failed to link adapter %s: %v", adapterID, err)
		}
	}

	if err := k.repository.Delete(device.ID); err != nil {
		return err
	}
	k.deleteMutex(device.ID)

	log.Printf("[Kernel] Device unregistered: %s (ID: %s)", device.Name, device.ID)

	return nil
}

func (ds *Kernel) GetDevice(deviceID string) (*types.Device, error) {
	return ds.repository.FindByID(deviceID)
}

func (ds *Kernel) ListDevices() ([]*types.Device, error) {
	return ds.repository.FindAll()
}

// --- Linking Logic ---

func (k *Kernel) LinkDeviceToAdapter(deviceID, adapterID string) error {
	device, err := k.repository.FindByID(deviceID)
	if err != nil || device == nil {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	adapter, err := k.pluginManager.GetPluginById(plugin.PluginAdapter, adapterID)
	if err != nil {
		return fmt.Errorf("adapter not found: %s", adapterID)
	}

	if err := k.repository.LinkAdapter(deviceID, adapterID); err != nil {
		return err
	}

	k.eventBus.Publish(events.Event{
		Type:    events.RegisterDeviceForAdapter(adapter.ID),
		Payload: device,
	})
	return nil
}

func (k *Kernel) UnlinkDeviceFromAdapter(deviceID, adapterID string) error {
	device, err := k.repository.FindByID(deviceID)
	if err != nil || device == nil {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	adapter, err := k.pluginManager.GetPluginById(plugin.PluginAdapter, adapterID)
	if err != nil {
		return fmt.Errorf("adapter not found: %s", adapterID)
	}

	if err := k.repository.UnlinkAdapter(deviceID, adapterID); err != nil {
		return err
	}

	k.eventBus.Publish(events.Event{
		Type:    events.UnregisterDeviceForAdapter(adapter.ID),
		Payload: device,
	})
	return nil
}

func (k *Kernel) ListPlugins() []*plugin.Plugin {
	plugins := k.pluginManager.GetPlugins()

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].ID < plugins[j].ID
	})
	return plugins
}

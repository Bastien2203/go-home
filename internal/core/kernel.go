package core

import (
	"fmt"
	"os/exec"

	"log"
	"sort"
	"sync"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/plugin"
	"github.com/Bastien2203/go-home/shared/types"
)

type Kernel struct {
	eventBus      *events.EventBus
	repository    DeviceRepository
	mu            map[string]*sync.Mutex
	muLock        sync.Mutex
	pluginManager *PluginManager
	processes     map[string]*exec.Cmd
}

func NewKernel(eventBus *events.EventBus, repository DeviceRepository) (*Kernel, error) {
	pluginManager, err := NewPluginManager(eventBus)
	if err != nil {
		return nil, err
	}

	kernel := &Kernel{
		eventBus:      eventBus,
		repository:    repository,
		mu:            make(map[string]*sync.Mutex),
		processes:     make(map[string]*exec.Cmd),
		pluginManager: pluginManager,
	}

	if err := events.Subscribe(eventBus, events.ParsedDataReceived, kernel.handleStateUpdate); err != nil {
		return nil, err
	}

	return kernel, nil
}

func (k *Kernel) handleStateUpdate(parsedData types.ParsedData) {
	device, err := k.repository.FindByAddress(parsedData.Address, parsedData.AddressType)
	if err != nil || device == nil {
		return
	}

	device.LastUpdated = parsedData.Timestamp

	mu := k.getMutex(device.ID)
	mu.Lock()
	for _, c := range parsedData.Data {
		device.Capabilities[c.Name] = c
	}
	// TODO : Save updated device (but its too much intensive to do it for each data update)
	k.repository.Save(device)
	mu.Unlock()
	fmt.Printf("Capabilities updated for device %s: %+v\n", device.Name, device.Capabilities)

	for _, adapterID := range device.AdapterIDs {
		go func(adapterID string) {
			for _, c := range parsedData.Data {
				k.eventBus.Publish(events.Event{
					Type: events.UpdateDataForAdapter(adapterID),
					Payload: types.DeviceStateUpdate{
						DeviceID:       device.ID,
						DeviceName:     device.Name,
						CapabilityType: c.Name,
						Timestamp:      parsedData.Timestamp,
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

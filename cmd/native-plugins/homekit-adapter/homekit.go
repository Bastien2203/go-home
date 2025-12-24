package main

import (
	"log"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
)

type HomekitAdapter struct {
	server  *HomekitServer
	manager *HomekitManager
}

func NewHomeKitAdapter(eventBus *events.EventBus, onStateChange func(state types.State), homekitDataDir string) (*HomekitAdapter, error) {
	server := NewHomekitServer(onStateChange, homekitDataDir)
	a := &HomekitAdapter{
		server:  server,
		manager: NewHomekitManager(server),
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
	log.Println("Server started, configuration code is '0010-2003'")
	return h.server.Start()
}

func (h *HomekitAdapter) Stop() error {
	return h.server.Stop()
}

func (h *HomekitAdapter) onDeviceData(data types.DeviceStateUpdate) {
	mapping, supported := CapabilityRegistry[data.CapabilityType]
	if !supported {
		log.Printf("capability %s not supported", data.CapabilityType)
		return
	}

	if !h.manager.AccessoryExists(data.DeviceID) {
		h.manager.CreateAccessory(data.DeviceName, data.DeviceID, mapping.Type)
		h.manager.CreateService(data.DeviceID, mapping.NewService())
	}

	svc := h.manager.GetService(data.DeviceID, mapping.ServiceType)

	if svc == nil {
		svc = h.manager.CreateService(data.DeviceID, mapping.NewService())
		if svc == nil {
			log.Printf("failed to create service capability %s for device %s", data.CapabilityType, data.DeviceID)
		}
	}

	if !h.manager.UdateCharacteristic(svc.C(mapping.CharType), mapping.ValueConverter(data.Value)) {
		log.Printf("failed to update charachteristics for device %s", data.DeviceID)
	}
}

func (h *HomekitAdapter) onDeviceRegistered(dev types.Device) {
	// Device auto register when data is received
}

func (h *HomekitAdapter) onDeviceUnregistered(dev types.Device) {
	h.manager.RemoveAccessory(dev.ID)
}

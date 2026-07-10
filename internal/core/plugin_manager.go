package core

import (
	"fmt"

	"log"
	"sync"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/plugin"
	"github.com/Bastien2203/go-home/utils"
)

type PluginManager struct {
	eventBus *events.EventBus
	plugins  map[plugin.PluginType]map[string]*plugin.Plugin
	mu       sync.Mutex

	pendingRequests map[string]chan error
	muRequests      sync.Mutex
}

const TimeoutDuration = 30 * time.Second

func NewPluginManager(eventBus *events.EventBus) (*PluginManager, error) {
	manager := &PluginManager{
		eventBus:        eventBus,
		plugins:         make(map[plugin.PluginType]map[string]*plugin.Plugin),
		pendingRequests: make(map[string]chan error),
	}

	if err := manager.subscribeToEvents(); err != nil {
		return nil, err
	}

	manager.eventBus.Publish(events.Event{
		Type:    events.CoreDiscovery,
		Payload: nil,
	})

	return manager, nil
}

func (m *PluginManager) subscribeToEvents() error {
	if err := events.Subscribe(m.eventBus, events.PluginConnected, m.onPluginConnected); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.PluginDisconnected, m.onPluginDisconnected); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.PluginStateChanged, m.onPluginStateChanged); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.PluginCommandResponse, m.onCommandResponse); err != nil {
		return err
	}

	// subscrive to other events here ...

	return nil
}

func (m *PluginManager) onCommandResponse(r plugin.CommandResponse) {
	m.muRequests.Lock()
	defer m.muRequests.Unlock()
	ch, exists := m.pendingRequests[r.RequestID]
	if exists {
		if !r.Success {
			ch <- fmt.Errorf("plugin error: %s", r.Error)
		} else {
			ch <- nil
		}
		delete(m.pendingRequests, r.RequestID)
	}
}

func (m *PluginManager) sendCommandAndWait(eventType events.EventType, pluginName string) error {
	reqID := fmt.Sprintf("%d-%s-%s", time.Now().UnixNano(), pluginName, eventType)
	ch := make(chan error, 1)

	m.muRequests.Lock()
	m.pendingRequests[reqID] = ch
	m.muRequests.Unlock()

	m.eventBus.Publish(events.Event{
		Type: eventType,
		Payload: plugin.CommandRequest{
			RequestID: reqID,
		},
	})

	select {
	case err := <-ch:
		return err
	case <-time.After(TimeoutDuration):
		m.muRequests.Lock()
		delete(m.pendingRequests, reqID)
		m.muRequests.Unlock()
		close(ch)
		return fmt.Errorf("timeout waiting for plugin %s", pluginName)
	}
}

func (m *PluginManager) onPluginConnected(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.plugins[p.Type][p.ID]; ok {
		log.Printf("[PluginManager] plugin with ID:%s already exists", p.ID)
		return
	}

	if _, ok := m.plugins[p.Type]; !ok {
		m.plugins[p.Type] = make(map[string]*plugin.Plugin)
	}

	m.plugins[p.Type][p.ID] = &p
}

func (m *PluginManager) onPluginDisconnected(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.plugins[p.Type][p.ID]; !ok {
		log.Printf("[PluginManager] plugin with ID:%s doesnt exists", p.ID)
		return
	}

	delete(m.plugins[p.Type], p.ID)
}

func (m *PluginManager) onPluginStateChanged(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.plugins[p.Type][p.ID]; !ok {
		log.Printf("[PluginManager] plugin with ID:%s do not exists", p.ID)
		return
	}

	m.plugins[p.Type][p.ID] = &p
}

func (m *PluginManager) GetPluginsByType(t plugin.PluginType) []*plugin.Plugin {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugins, ok := m.plugins[t]
	if !ok {
		return []*plugin.Plugin{}
	}

	return utils.Values(plugins)
}

func (m *PluginManager) GetPluginById(t plugin.PluginType, id string) (*plugin.Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, ok := m.plugins[t][id]
	if !ok {
		return nil, fmt.Errorf("plugin with id : %s and type : %s not found", id, t)
	}

	return plugin, nil
}

func (m *PluginManager) GetPlugins() []*plugin.Plugin {
	m.mu.Lock()
	defer m.mu.Unlock()

	total := 0
	for _, group := range m.plugins {
		total += len(group)
	}

	plugins := make([]*plugin.Plugin, 0, total)
	for _, t := range plugin.PluginTypes {
		plugins = append(plugins, utils.Values(m.plugins[t])...)
	}
	return plugins
}

func (m *PluginManager) StopPlugin(p *plugin.Plugin) error {
	return m.sendCommandAndWait(events.PluginStop(p.ID), p.Name)
}

func (m *PluginManager) StartPlugin(p *plugin.Plugin) error {
	return m.sendCommandAndWait(events.PluginStart(p.ID), p.Name)
}

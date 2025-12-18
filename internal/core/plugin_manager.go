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
	eventBus    *events.EventBus
	plugins     map[plugin.PluginType]map[string]*plugin.Plugin
	mu          sync.Mutex
	ack         map[string]chan struct{}
	negativeAck map[string]chan struct{}
}

const TimeoutDuration = 5 * time.Second

func NewPluginManager(eventBus *events.EventBus) (*PluginManager, error) {
	manager := &PluginManager{
		eventBus:    eventBus,
		plugins:     make(map[plugin.PluginType]map[string]*plugin.Plugin),
		ack:         make(map[string]chan struct{}),
		negativeAck: make(map[string]chan struct{}),
	}

	if err := manager.subscribeToEvents(); err != nil {
		return nil, err
	}

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

	if err := events.Subscribe(m.eventBus, events.PluginAck, m.onPluginAck); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.PluginNegativeAck, m.onPluginNegativeAck); err != nil {
		return err
	}

	// subscrive to other events here ...

	return nil
}

func (m *PluginManager) onPluginAck(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ack[p.ID] <- struct{}{}
}

func (m *PluginManager) onPluginNegativeAck(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.negativeAck[p.ID] <- struct{}{}
}

func (m *PluginManager) onPluginConnected(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.plugins[p.Type][p.ID]; ok {
		log.Printf("[PluginManager] plugin with ID:%s already exists", p.ID)
		return
	}

	if _, ok := m.ack[p.ID]; ok {
		log.Printf("[PluginManager] plugin with ID:%s already exists in ack map", p.ID)
		return
	}

	if _, ok := m.negativeAck[p.ID]; ok {
		log.Printf("[PluginManager] plugin with ID:%s already exists in negativeAck map", p.ID)
		return
	}
	if _, ok := m.plugins[p.Type]; !ok {
		m.plugins[p.Type] = make(map[string]*plugin.Plugin)
	}

	m.negativeAck[p.ID] = make(chan struct{}, 1)
	m.ack[p.ID] = make(chan struct{}, 1)
	m.plugins[p.Type][p.ID] = &p
}

func (m *PluginManager) onPluginDisconnected(p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.plugins[p.Type][p.ID]; !ok {
		log.Printf("[PluginManager] plugin with ID:%s doesnt exists", p.ID)
		return
	}

	if _, ok := m.ack[p.ID]; !ok {
		log.Printf("[PluginManager] plugin with ID:%s doesnt exists in ack map", p.ID)
		return
	}

	if _, ok := m.negativeAck[p.ID]; !ok {
		log.Printf("[PluginManager] plugin with ID:%s doesnt exists in negativeAck map", p.ID)
		return
	}

	delete(m.negativeAck, p.ID)
	delete(m.ack, p.ID)
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
	m.eventBus.Publish(events.Event{
		Type:    events.PluginStop(p.ID),
		Payload: []any{},
	})

	select {
	case <-m.ack[p.ID]:
		return nil
	case <-m.negativeAck[p.ID]:
		return fmt.Errorf("error stopping %s %s", p.Type, p.Name)
	case <-time.After(TimeoutDuration):
		return fmt.Errorf("timeout stopping %s %s", p.Type, p.Name)
	}
}

func (m *PluginManager) StartPlugin(p *plugin.Plugin) error {
	m.mu.Lock()
	ack := m.ack[p.ID]
	nack := m.negativeAck[p.ID]
	m.mu.Unlock()
	m.eventBus.Publish(events.Event{
		Type:    events.PluginStart(p.ID),
		Payload: []any{},
	})

	select {
	case <-ack:
		return nil
	case <-nack:
		return fmt.Errorf("error starting %s %s", p.Type, p.Name)
	case <-time.After(TimeoutDuration):
		return fmt.Errorf("timeout starting %s %s", p.Type, p.Name)
	}
}

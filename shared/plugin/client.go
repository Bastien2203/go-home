package plugin

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
)

type PluginClient struct {
	pluginInstance *Plugin
	eventBus       *events.EventBus
	onStart        func() error
	onStop         func() error
}

func NewPluginClient(instance *Plugin, eventBus *events.EventBus) *PluginClient {
	client := &PluginClient{
		eventBus:       eventBus,
		pluginInstance: instance,
	}

	return client
}

func (m *PluginClient) subscribeToEvents() error {
	if err := events.Subscribe(m.eventBus, events.PluginStop(m.pluginInstance.ID), m.onPluginStop); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.PluginStart(m.pluginInstance.ID), m.onPluginStart); err != nil {
		return err
	}

	if err := events.Subscribe(m.eventBus, events.CoreDiscovery, m.onCoreDiscovery); err != nil {
		return err
	}

	return nil
}

func (c *PluginClient) RunPlugin(onStart func() error, onStop func() error) {
	c.onStart = onStart
	c.onStop = onStop
	if err := c.subscribeToEvents(); err != nil {
		log.Fatalf("error while subscribing to events : %v", err)
	}

	c.announcePresence()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	c.eventBus.Publish(events.Event{
		Type:    events.PluginDisconnected,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) announcePresence() {
	c.eventBus.Publish(events.Event{
		Type:    events.PluginConnected,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) onCoreDiscovery(_ any) {
	log.Printf("Core discovery received, re-announcing presence")
	c.announcePresence()
}

func (c *PluginClient) replyCommand(reqID string, err error) {
	resp := CommandResponse{
		RequestID: reqID,
		PluginID:  c.pluginInstance.ID,
		Success:   err == nil,
	}
	if err != nil {
		resp.Error = err.Error()
	}

	c.eventBus.Publish(events.Event{
		Type:    events.PluginCommandResponse,
		Payload: resp,
	})
}

func (c *PluginClient) onPluginStop(req CommandRequest) {
	err := c.onStop()
	if err != nil {
		log.Printf("error on plugin stop : %v", err)
	}
	c.replyCommand(req.RequestID, err)
}

func (c *PluginClient) onPluginStart(req CommandRequest) {
	err := c.onStart()
	if err != nil {
		log.Printf("error on plugin start : %v", err)
	}
	c.replyCommand(req.RequestID, err)
}

func (c *PluginClient) EmitNewState(s types.State) {
	c.pluginInstance.State = s
	c.eventBus.Publish(events.Event{
		Type:    events.PluginStateChanged,
		Payload: c.pluginInstance,
	})
}

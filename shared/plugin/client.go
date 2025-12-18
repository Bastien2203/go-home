package plugin

import (
	"gohome/shared/events"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gohome/shared/types"
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

	return nil
}

func (c *PluginClient) RunPlugin(onStart func() error, onStop func() error) {
	c.onStart = onStart
	c.onStop = onStop
	if err := c.subscribeToEvents(); err != nil {
		log.Fatalf("error while subscribing to events : %v", err)
	}
	c.eventBus.Publish(events.Event{
		Type:    events.PluginConnected,
		Payload: c.pluginInstance,
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	c.eventBus.Publish(events.Event{
		Type:    events.PluginDisconnected,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) EmitNewState(s types.State) {
	c.pluginInstance.State = s
	c.eventBus.Publish(events.Event{
		Type:    events.PluginStateChanged,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) ack() {
	c.eventBus.Publish(events.Event{
		Type:    events.PluginAck,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) negativeAck() {
	c.eventBus.Publish(events.Event{
		Type:    events.PluginNegativeAck,
		Payload: c.pluginInstance,
	})
}

func (c *PluginClient) onPluginStop(_ []any) {
	if err := c.onStop(); err != nil {
		log.Printf("error on plugin stop : %v", err)
		c.negativeAck()
		return
	}
	c.ack()
}

func (c *PluginClient) onPluginStart(_ []any) {
	if err := c.onStart(); err != nil {
		log.Printf("error on plugin start : %v", err)
		c.negativeAck()
		return
	}
	c.ack()
}

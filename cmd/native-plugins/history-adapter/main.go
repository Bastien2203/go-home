package main

import (
	"context"
	"gohome/shared/config"
	"gohome/shared/events"
	"gohome/shared/plugin"
	"gohome/shared/types"
	"log"
)

var p = &plugin.Plugin{
	ID:    "history-adapter",
	Name:  "History",
	Type:  plugin.PluginAdapter,
	State: types.StateStopped,
	Widgets: map[string]*plugin.Widget{
		"history-adapter-widget": {
			ID:   "history-adapter-widget",
			Type: plugin.TypeLineChart,
			Name: "Data history",
			Config: map[string]any{
				"dataUrl": ":8888/api/history/device/{deviceId}/capabilities/{capabilityType}",
			},
			MountPoint: plugin.CapabilityWidget,
		},
	},
}

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnv(ctx)

	eventBus, err := events.NewEventBus(cfg.BrokerUrl, p.ID)
	if err != nil {
		log.Fatalf("Error setting up event bus : %v", err)
	}

	client := plugin.NewPluginClient(p, eventBus)
	adapter, err := NewHistoryAdapter(eventBus, client.EmitNewState, 8888)
	if err != nil {
		log.Fatalf("Error creating homekit adapter : %v", err)
	}
	client.RunPlugin(adapter.Start, adapter.Stop)
}

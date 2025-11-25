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
	ID:    "homekit-adapter",
	Name:  "Homekit",
	Type:  plugin.PluginAdapter,
	State: types.StateStopped,
}

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnv(ctx)

	eventBus, err := events.NewEventBus(cfg.BrokerUrl, p.ID)
	if err != nil {
		log.Fatalf("Error setting up event bus : %v", err)
	}

	client := plugin.NewPluginClient(p, eventBus)
	adapter, err := NewHomeKitAdapter(eventBus, client.EmitNewState, "./homekit_data")
	if err != nil {
		log.Fatalf("Error creating homekit adapter : %v", err)
	}
	client.RunPlugin(adapter.Start, adapter.Stop)
}

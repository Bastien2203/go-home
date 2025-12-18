package main

import (
	"context"

	"log"

	"github.com/Bastien2203/go-home/shared/config"
	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/plugin"
	"github.com/Bastien2203/go-home/shared/types"
)

var p = &plugin.Plugin{
	ID:      "homekit-adapter",
	Name:    "Homekit",
	Type:    plugin.PluginAdapter,
	State:   types.StateStopped,
	Widgets: map[string]*plugin.Widget{},
}

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnvPlugin(ctx)

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

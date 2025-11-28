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
	ID:      "http-scanner",
	Name:    "Http Scanner",
	Type:    plugin.PluginScanner,
	State:   types.StateStopped,
	Widgets: map[string]*plugin.Widget{},
}

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnv(ctx)

	eventBus, err := events.NewEventBus(cfg.BrokerUrl, p.ID)
	if err != nil {
		log.Fatalf("Error setting up event bus : %v", err)
	}

	client := plugin.NewPluginClient(p, eventBus)
	scanner := NewHTTPScanner(eventBus, 8889, client.EmitNewState)
	client.RunPlugin(scanner.Start, scanner.Stop)
}

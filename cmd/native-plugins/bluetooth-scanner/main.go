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
	ID:    "bluetooth-scanner",
	Name:  "Bluetooth Scanner",
	Type:  plugin.PluginScanner,
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
	scanner := NewBluetoothScanner(eventBus, client.EmitNewState)
	client.RunPlugin(scanner.Start, scanner.Stop)
}

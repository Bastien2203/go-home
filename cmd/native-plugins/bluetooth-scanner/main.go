package main

import (
	"context"
	"log"
	"os"

	"github.com/Bastien2203/go-home/shared/config"
	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/plugin"
	"github.com/Bastien2203/go-home/shared/types"
)

var p = &plugin.Plugin{
	ID:    "bluetooth-scanner",
	Name:  "Bluetooth Scanner",
	Type:  plugin.PluginScanner,
	State: types.StateStopped,
}

func main() {
	ctx := context.Background()
	cfg := config.LoadFromEnvPlugin(ctx)

	eventBus, err := events.NewEventBus(cfg.BrokerUrl, p.ID)
	if err != nil {
		log.Fatalf("Error setting up event bus : %v", err)
	}

	client := plugin.NewPluginClient(p, eventBus)
	scanner := NewBluetoothScanner(eventBus, client.EmitNewState)
	if os.Getenv("QUICK_STARTUP") == "true" {
		log.Printf("Quick startup enabled")
		scanner.Start()
	}
	client.RunPlugin(scanner.Start, scanner.Stop)
}

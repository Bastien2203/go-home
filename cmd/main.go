package main

import (
	"context"
	"fmt"
	"gohome/internal/adapters"
	"gohome/internal/core"
	"gohome/internal/events"
	"gohome/internal/protocols"
	"gohome/internal/repository"
	"gohome/internal/scanners"
	"gohome/internal/server"
	"gohome/internal/websockets"
	"log"
	"os"
	"os/signal"
)

func setupAdapters(kernel *core.Kernel) {
	printerAdapter := adapters.NewPrinterAdapter(kernel)
	if err := kernel.RegisterAdapter(printerAdapter); err != nil {
		log.Fatal(err)
	}

	homekitAdapter := adapters.NewHomeKitAdapter(kernel)
	if err := kernel.RegisterAdapter(homekitAdapter); err != nil {
		log.Fatal(err)
	}
}

func setupProtocols(kernel *core.Kernel) {
	dummyParser := protocols.NewHttpParser()
	bthomeParser := protocols.NewBthomeParser()
	kernel.RegisterProtocol(dummyParser)
	kernel.RegisterProtocol(bthomeParser)
}

func setupScanners(ctx context.Context, eventBus *events.EventBus, kernel *core.Kernel) {
	bluetoothScanner := scanners.NewBluetoothScanner(eventBus)
	if err := kernel.RegisterScanner(bluetoothScanner, ctx); err != nil {
		log.Fatal(err)
	}

	httpScanner := scanners.NewHTTPScanner(eventBus, 8888)
	if err := kernel.RegisterScanner(httpScanner, ctx); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventBus := events.NewEventBus()
	db, err := repository.SetupSQLiteDB("./gohome.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deviceRepo, err := repository.NewSQLiteDeviceRepository(db)
	if err != nil {
		log.Fatalf("Error init sqlite repo: %v", err)
	}
	kernel := core.NewKernel(eventBus, deviceRepo)

	setupAdapters(kernel)
	setupProtocols(kernel)
	setupScanners(ctx, eventBus, kernel)

	defer kernel.Stop()

	wsHub := websockets.NewHub()
	go wsHub.Run()
	apiServer := server.NewServer(kernel, 8080, wsHub)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	eventBus.Subscribe(events.BluetoothDeviceFound, func(event events.Event) {
		wsHub.Broadcast(websockets.TopicBluetoothDevice, event.Payload)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nShutting down...")
}

//addr := "d69eee9c-8848-5c4c-3c41-d5b88b929976"
//name := "Temperature Sensor"

// Create and register a thermometer device ------------------------------------------------
// thermometer := &core.Device{
// 	ID:           "d69eee9c-8848-5c4c-3c41-d5b88b929976",
// 	Name:         "Living Room Thermometer",
// 	Type:         core.TemperatureSensor,
// 	Protocol:     bthomeParser.Name(),
// 	AdapterIDs:   []string{homekitAdapter.ID()},
// 	CreatedAt:    time.Now(),
// 	Capabilities: map[core.CapabilityType]*core.Capability{},
// }

// if err := kernel.RegisterDevice(thermometer); err != nil {
// 	log.Fatal(err)
// }

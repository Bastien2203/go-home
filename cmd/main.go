package main

import (
	"context"
	"fmt"
	"gohome/internal/adapters"
	"gohome/internal/core"
	"gohome/internal/events"
	"gohome/internal/protocols"
	"gohome/internal/scanners"
	"gohome/internal/server"
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
	dummyParser := protocols.NewDummyParser()
	bthomeParser := protocols.NewBthomeParser()
	kernel.RegisterProtocol(dummyParser)
	kernel.RegisterProtocol(bthomeParser)
}

func setupScanners(ctx context.Context, eventBus *events.EventBus, kernel *core.Kernel) {
	bluetoothScanner := scanners.NewBluetoothScanner(eventBus)
	if err := kernel.RegisterScanner(bluetoothScanner, ctx); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventBus := events.NewEventBus()
	repository := core.NewInMemoryDeviceRepository()
	kernel := core.NewKernel(eventBus, repository)

	setupAdapters(kernel)
	setupProtocols(kernel)
	setupScanners(ctx, eventBus, kernel)

	defer kernel.Stop()

	apiServer := server.NewServer(kernel, 8080)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

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

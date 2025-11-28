package main

import (
	"context"
	"fmt"
	"gohome/internal/core"
	"gohome/internal/repository"

	"gohome/shared/config"
	"gohome/shared/events"

	"gohome/internal/protocols"

	"gohome/internal/server"
	"gohome/internal/websockets"
	"log"
	"os"
	"os/signal"
)

func setupProtocols(kernel *core.Kernel) {
	dummyParser := protocols.NewHttpParser()
	bthomeParser := protocols.NewBthomeParser()
	kernel.RegisterProtocol(dummyParser)
	kernel.RegisterProtocol(bthomeParser)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadFromEnv(ctx)

	eventBus, err := events.NewEventBus(cfg.BrokerUrl, "gohome-core")
	if err != nil {
		log.Fatalf("Failed to start event bus : %v", err)
	}

	defer eventBus.Close()

	db, err := repository.SetupSQLiteDB(cfg.SqliteDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deviceRepo, err := repository.NewDeviceRepository(db)
	if err != nil {
		log.Fatalf("Error init sqlite device repo: %v", err)
	}

	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("Error init sqlite users repo: %v", err)
	}
	kernel, err := core.NewKernel(eventBus, deviceRepo)
	if err != nil {
		log.Fatalf("Failed to create kernel: %v", err)
	}

	// TODO: handle protocols as plugin
	setupProtocols(kernel)

	wsHub := websockets.NewHub()
	go wsHub.Run()
	apiServer := server.NewServer(kernel, cfg.ApiPort, cfg.SessionSecret, cfg.AppEnv, wsHub, userRepo)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	if err := events.Subscribe(eventBus, events.BluetoothDeviceFound, func(payload any) {
		wsHub.Broadcast(websockets.TopicBluetoothDevice, payload)
	}); err != nil {
		log.Fatalf("Failed to subscribe to event topic bluetooth discovery: %v", err)
	}

	kernel.LoadPlugins(cfg.PluginFolderPath)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Stopping scanners...")
	kernel.StopScanners()

	fmt.Println("Stopping adapters...")
	kernel.StopAdapters()

	fmt.Println("Unloading plugins...")
	kernel.UnloadPlugins()

	fmt.Println("\nShutting down...")
}

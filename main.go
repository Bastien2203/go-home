package main

import (
	"context"
	"fmt"

	"log"
	"os"
	"os/signal"

	"github.com/Bastien2203/go-home/internal/core"
	"github.com/Bastien2203/go-home/internal/repository"
	"github.com/Bastien2203/go-home/internal/server"
	"github.com/Bastien2203/go-home/internal/websockets"
	"github.com/Bastien2203/go-home/shared/config"
	"github.com/Bastien2203/go-home/shared/events"
)

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("Stopping scanners...")
	kernel.StopScanners()

	fmt.Println("Stopping adapters...")
	kernel.StopAdapters()

	fmt.Println("\nShutting down...")
}

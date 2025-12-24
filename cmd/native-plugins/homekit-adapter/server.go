package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/Bastien2203/go-home/shared/types"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

type HomekitServer struct {
	onStateChange  func(state types.State)
	server         *hap.Server
	serverStop     context.CancelFunc
	mu             sync.Mutex
	restartTimer   *time.Timer
	homekitDataDir string
	bridge         *accessory.Bridge
}

func NewHomekitServer(onStateChange func(state types.State), homekitDataDir string) *HomekitServer {
	info := accessory.Info{
		Name:         "DEV GoHome Hub",
		SerialNumber: "GOHOME-HUB-001",
		Manufacturer: "GoHome",
		Model:        "Bridge v2",
		Firmware:     "1.0.0",
	}
	bridge := accessory.NewBridge(info)
	return &HomekitServer{
		onStateChange:  onStateChange,
		homekitDataDir: homekitDataDir,
		bridge:         bridge,
	}
}

func (s *HomekitServer) Start() error {
	// Cannot start without at least one device
	return nil
}

func (s *HomekitServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.serverStop != nil {
		s.serverStop()
		s.onStateChange(types.StateStopped)
		s.serverStop = nil
	}
	log.Println("[HomeKit] Adapter stopped.")
	return nil
}

func (s *HomekitServer) ScheduleReload(accessories map[string]*accessory.A) {
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.restartTimer != nil {
			s.restartTimer.Stop()
		}

		log.Println("[HomeKit] Structure change detected. Scheduling reload in 2s...")

		s.restartTimer = time.AfterFunc(2*time.Second, func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			log.Println("[HomeKit] Executing scheduled server reload...")
			if err := s.reloadServer(accessories); err != nil {
				log.Printf("[HomeKit] Error reloading server: %v", err)
			}
			s.restartTimer = nil
		})
	}()
}

func (s *HomekitServer) reloadServer(accessories map[string]*accessory.A) error {

	s.onStateChange(types.StateStopped)
	if s.serverStop != nil {
		s.serverStop()
		s.serverStop = nil
		time.Sleep(100 * time.Millisecond)
	}

	var accs []*accessory.A
	for _, a := range accessories {
		accs = append(accs, a)
	}

	if len(accs) == 0 {
		log.Println("[HomeKit] No accessories to publish yet.")
		return nil
	}

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].Id < accs[j].Id
	})

	fs := hap.NewFsStore(s.homekitDataDir)

	server, err := hap.NewServer(fs, s.bridge.A, accs...)
	if err != nil {
		return fmt.Errorf("failed to create hap server: %w", err)
	}
	server.Ifaces = []string{os.Getenv("INTERNET_INTERFACE")}
	s.server = server

	ctx, cancel := context.WithCancel(context.Background())
	s.serverStop = cancel

	go func() {
		log.Printf("[HomeKit] Server running with %d accessories", len(accs))
		if err := server.ListenAndServe(ctx); err != nil && err != context.Canceled {
			log.Printf("[HomeKit] Server stopped with error: %v", err)
		}
	}()

	s.onStateChange(types.StateRunning)
	return nil
}

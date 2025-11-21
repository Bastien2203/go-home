package scanners

import (
	"context"
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"gohome/internal/events"
	"log"
	"net/http"
	"time"
)

type HTTPScanner struct {
	eventBus *events.EventBus
	address  string // ex: ":8080"
	server   *http.Server
	started  bool
}

func NewHTTPScanner(eventBus *events.EventBus, port int) core.Scanner {
	addr := fmt.Sprintf(":%d", port)
	return &HTTPScanner{
		eventBus: eventBus,
		address:  addr,
		started:  false,
	}
}

func (s *HTTPScanner) ID() string {
	return "http_scanner"
}

func (s *HTTPScanner) Name() string {
	return "HTTP Scanner"
}

func (s *HTTPScanner) State() core.State {
	switch s.started {
	case true:
		return core.StateRunning
	default:
		return core.StateStopped
	}
}

func (s *HTTPScanner) Start(ctx context.Context) error {
	if s.started {
		return fmt.Errorf("http listener already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/http-scanner/{addr}", s.handleRequest)

	s.server = &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	go func() {
		log.Printf("[HTTP Listener] Listening on %s/api/http-scanner/{addr}", s.address)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[HTTP Listener] Server error: %v", err)
			s.started = false
		}
	}()

	s.started = true
	return nil
}

func (s *HTTPScanner) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addr := r.PathValue("addr")

	var payloads []*core.Capability

	if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
		log.Printf("[HTTP Listener] Failed to decode JSON from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("[HTTP Listener] Received data from %s: %d metrics", r.RemoteAddr, len(payloads))

	if s.eventBus != nil {
		s.eventBus.Publish(events.Event{
			Type: events.RawDataReceived,
			Payload: &core.RawData{
				Address:     addr,
				Data:        payloads,
				Timestamp:   time.Now(),
				AddressType: core.BasicAddress,
			},
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *HTTPScanner) Stop() error {
	if !s.started {
		return nil
	}

	log.Println("[HTTP Listener] Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("[HTTP Listener] Shutdown error: %v", err)
		return err
	}

	s.started = false
	log.Println("[HTTP Listener] Stopped")
	return nil
}

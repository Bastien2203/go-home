package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/Bastien2203/go-home/shared/events"
	"github.com/Bastien2203/go-home/shared/types"
)

type HTTPScanner struct {
	eventBus      *events.EventBus
	address       string // ex: ":8080"
	server        *http.Server
	started       bool
	onStateChange func(state types.State)
}

func NewHTTPScanner(eventBus *events.EventBus, port int, onStateChange func(state types.State)) *HTTPScanner {
	addr := fmt.Sprintf(":%d", port)
	return &HTTPScanner{
		eventBus:      eventBus,
		address:       addr,
		started:       false,
		onStateChange: onStateChange,
	}
}

func (s *HTTPScanner) Start() error {
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
			s.onStateChange(types.StateStopped)
			s.started = false
		}
	}()

	s.onStateChange(types.StateRunning)
	s.started = true
	return nil
}

func (s *HTTPScanner) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addr := r.PathValue("addr")

	var payloads []*types.Capability

	if err := json.NewDecoder(r.Body).Decode(&payloads); err != nil {
		log.Printf("[HTTP Listener] Failed to decode JSON from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("[HTTP Listener] Received data from %s: %d metrics", r.RemoteAddr, len(payloads))

	if s.eventBus != nil {
		bytes, err := json.Marshal(payloads)
		if err != nil {
			log.Printf("[HTTP Listener] Failed to marshal  %v", err)
		}
		s.eventBus.Publish(events.Event{
			Type: events.RawDataReceived,
			Payload: &types.RawData{
				Address:     addr,
				Data:        bytes,
				Timestamp:   time.Now(),
				AddressType: types.BasicAddress,
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
	s.onStateChange(types.StateStopped)
	log.Println("[HTTP Listener] Stopped")
	return nil
}

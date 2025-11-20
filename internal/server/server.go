package server

import (
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"log"
	"net/http"
)

type Server struct {
	kernel *core.Kernel
	addr   string
}

func NewServer(kernel *core.Kernel, port int) *Server {
	return &Server{
		kernel: kernel,
		addr:   fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	handler := s.corsMiddleware(mux)

	// --- Routes ---

	mux.HandleFunc("GET /api/adapters", s.handleListAdapters)
	mux.HandleFunc("GET /api/scanners", s.handleListScanners)
	mux.HandleFunc("GET /api/protocols", s.handleListProtocols)
	mux.HandleFunc("GET /api/devices", s.handleListDevices)

	mux.HandleFunc("POST /api/devices", s.handleCreateDevice)

	mux.HandleFunc("POST /api/devices/{id}/adapters/{adapterId}", s.handleLinkAdapter)
	mux.HandleFunc("DELETE /api/devices/{id}/adapters/{adapterId}", s.handleUnlinkAdapter)

	server := &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	log.Printf("[Server] API listening on http://localhost%s", s.addr)
	return server.ListenAndServe()
}

// --- Handlers ---

func (s *Server) handleListAdapters(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListAdapters())
}

func (s *Server) handleListScanners(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListScanners())
}

func (s *Server) handleListProtocols(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListProtocols())
}

func (s *Server) handleListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.kernel.ListDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if devices == nil {
		devices = []*core.Device{}
	}
	json.NewEncoder(w).Encode(devices)
}

func (s *Server) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	var req DeviceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	protocol, err := s.kernel.GetProtocol(req.Protocol)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unknown protocol: %s", req.Protocol), http.StatusBadRequest)
		return
	}

	dev := core.NewDevice(req.Address, req.Name, req.Protocol, req.AdapterIDs, protocol.AddressType())
	if err := s.kernel.RegisterDevice(dev); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register device: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dev)
}

func (s *Server) handleLinkAdapter(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	adapterID := r.PathValue("adapterId")

	if err := s.kernel.LinkDeviceToAdapter(deviceID, adapterID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "linked"}`))
}

func (s *Server) handleUnlinkAdapter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "Unlink adapter not implemented yet"}`))
}

// --- Middleware ---

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type DeviceCreateRequest struct {
	Address    string   `json:"address"`
	Name       string   `json:"name"`
	Protocol   string   `json:"protocol"`
	AdapterIDs []string `json:"adapter_ids"`
}

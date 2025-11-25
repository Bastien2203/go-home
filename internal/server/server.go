package server

import (
	"encoding/json"
	"fmt"
	"gohome/internal/core"
	"gohome/internal/websockets"
	"gohome/shared/types"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	kernel *core.Kernel
	addr   string
	wsHub  *websockets.Hub
}

func NewServer(kernel *core.Kernel, port int, wsHub *websockets.Hub) *Server {
	return &Server{
		kernel: kernel,
		addr:   fmt.Sprintf(":%d", port),
		wsHub:  wsHub,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticDir := "./dist"

	handler := s.corsMiddleware(mux)

	// --- Routes ---

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(s.wsHub, w, r)
	})

	mux.HandleFunc("GET /api/adapters", s.handleListAdapters)
	mux.HandleFunc("GET /api/scanners", s.handleListScanners)
	mux.HandleFunc("GET /api/protocols", s.handleListProtocols)
	mux.HandleFunc("GET /api/devices", s.handleListDevices)

	mux.HandleFunc("POST /api/devices", s.handleCreateDevice)
	mux.HandleFunc("DELETE /api/devices/{id}", s.handleDeleteDevice)

	mux.HandleFunc("POST /api/devices/{id}/adapters/{adapterId}", s.handleLinkAdapter)
	mux.HandleFunc("DELETE /api/devices/{id}/adapters/{adapterId}", s.handleUnlinkAdapter)

	mux.HandleFunc("POST /api/scanners/start/{scannerId}", s.handleStartScanner)
	mux.HandleFunc("POST /api/scanners/stop/{scannerId}", s.handleStopScanner)
	mux.HandleFunc("POST /api/adapters/start/{adapterId}", s.handleStartAdapter)
	mux.HandleFunc("POST /api/adapters/stop/{adapterId}", s.handleStopAdapter)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)

		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
	})

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
		devices = []*types.Device{}
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

	dev := types.NewDevice(req.Address, req.Name, req.Protocol, req.AdapterIDs, protocol.AddressType())
	if err := s.kernel.RegisterDevice(dev); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register device: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dev)
}

func (s *Server) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")

	if err := s.kernel.UnregisterDevice(deviceID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "deleted"}`))
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
	deviceID := r.PathValue("id")
	adapterID := r.PathValue("adapterId")

	if err := s.kernel.UnlinkDeviceFromAdapter(deviceID, adapterID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "unlinked"}`))
}

func (s *Server) handleStartScanner(w http.ResponseWriter, r *http.Request) {
	scannerID := r.PathValue("scannerId")

	if err := s.kernel.StartScanner(scannerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "started"}`))
}

func (s *Server) handleStopScanner(w http.ResponseWriter, r *http.Request) {
	scannerID := r.PathValue("scannerId")

	if err := s.kernel.StopScanner(scannerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "stopped"}`))
}

func (s *Server) handleStartAdapter(w http.ResponseWriter, r *http.Request) {
	adapterId := r.PathValue("adapterId")

	if err := s.kernel.StartAdapter(adapterId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "started"}`))
}

func (s *Server) handleStopAdapter(w http.ResponseWriter, r *http.Request) {
	adapterId := r.PathValue("adapterId")

	if err := s.kernel.StopAdapter(adapterId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "stopped"}`))
}

// --- Middleware ---

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if strings.HasPrefix(r.URL.Path, "/api") {
			w.Header().Set("Content-Type", "application/json")
		}

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

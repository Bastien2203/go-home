package routes

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/Bastien2203/go-home/internal/core"
	"github.com/Bastien2203/go-home/shared/types"
)

type DevicesRouter struct {
	kernel *core.Kernel
}

type DeviceCreateRequest struct {
	Address    string   `json:"address"`
	Name       string   `json:"name"`
	Protocol   string   `json:"protocol"`
	AdapterIDs []string `json:"adapter_ids"`
}

func NewDevicesRouter(kernel *core.Kernel, mux *http.ServeMux, middleware func(next http.Handler) http.Handler) *DevicesRouter {
	r := &DevicesRouter{
		kernel: kernel,
	}

	mux.Handle("GET /api/devices", middleware(http.HandlerFunc(r.handleListDevices)))
	mux.Handle("POST /api/devices", middleware(http.HandlerFunc(r.handleCreateDevice)))
	mux.Handle("DELETE /api/devices/{id}", middleware(http.HandlerFunc(r.handleDeleteDevice)))

	return r
}

func (s *DevicesRouter) handleListDevices(w http.ResponseWriter, r *http.Request) {
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

func (s *DevicesRouter) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
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

func (s *DevicesRouter) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")

	if err := s.kernel.UnregisterDevice(deviceID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "deleted"}`))
}

package routes

import (
	"encoding/json"

	"github.com/Bastien2203/go-home/internal/core"

	"net/http"
)

type AdaptersRouter struct {
	kernel *core.Kernel
}

func NewAdaptersRouter(kernel *core.Kernel, mux *http.ServeMux, middleware func(next http.Handler) http.Handler) *AdaptersRouter {
	r := &AdaptersRouter{
		kernel: kernel,
	}

	mux.Handle("GET /api/adapters", middleware(http.HandlerFunc(r.handleListAdapters)))
	mux.Handle("POST /api/devices/{id}/adapters/{adapterId}", middleware(http.HandlerFunc(r.handleLinkAdapter)))
	mux.Handle("DELETE /api/devices/{id}/adapters/{adapterId}", middleware(http.HandlerFunc(r.handleUnlinkAdapter)))
	mux.Handle("POST /api/adapters/start/{adapterId}", middleware(http.HandlerFunc(r.handleStartAdapter)))
	mux.Handle("POST /api/adapters/stop/{adapterId}", middleware(http.HandlerFunc(r.handleStopAdapter)))

	return r
}

func (s *AdaptersRouter) handleListAdapters(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListAdapters())
}

func (s *AdaptersRouter) handleLinkAdapter(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	adapterID := r.PathValue("adapterId")

	if err := s.kernel.LinkDeviceToAdapter(deviceID, adapterID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "linked"}`))
}

func (s *AdaptersRouter) handleUnlinkAdapter(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")
	adapterID := r.PathValue("adapterId")

	if err := s.kernel.UnlinkDeviceFromAdapter(deviceID, adapterID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "unlinked"}`))
}

func (s *AdaptersRouter) handleStartAdapter(w http.ResponseWriter, r *http.Request) {
	adapterId := r.PathValue("adapterId")

	if err := s.kernel.StartAdapter(adapterId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "started"}`))
}

func (s *AdaptersRouter) handleStopAdapter(w http.ResponseWriter, r *http.Request) {
	adapterId := r.PathValue("adapterId")

	if err := s.kernel.StopAdapter(adapterId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "stopped"}`))
}

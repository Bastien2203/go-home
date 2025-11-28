package routes

import (
	"encoding/json"
	"gohome/internal/core"
	"net/http"
)

type ScannersRouter struct {
	kernel *core.Kernel
}

func NewScannersRouter(kernel *core.Kernel, mux *http.ServeMux, middleware func(next http.Handler) http.Handler) *ScannersRouter {
	r := &ScannersRouter{
		kernel: kernel,
	}
	mux.Handle("GET /api/scanners", middleware(http.HandlerFunc(r.handleListScanners)))
	mux.Handle("POST /api/scanners/start/{scannerId}", middleware(http.HandlerFunc(r.handleStartScanner)))
	mux.Handle("POST /api/scanners/stop/{scannerId}", middleware(http.HandlerFunc(r.handleStopScanner)))
	return r
}

func (s *ScannersRouter) handleListScanners(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListScanners())
}

func (s *ScannersRouter) handleStartScanner(w http.ResponseWriter, r *http.Request) {
	scannerID := r.PathValue("scannerId")

	if err := s.kernel.StartScanner(scannerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "started"}`))
}

func (s *ScannersRouter) handleStopScanner(w http.ResponseWriter, r *http.Request) {
	scannerID := r.PathValue("scannerId")

	if err := s.kernel.StopScanner(scannerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "stopped"}`))
}

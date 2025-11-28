package routes

import (
	"encoding/json"
	"gohome/internal/core"
	"net/http"
)

type ProtocolsRouter struct {
	kernel *core.Kernel
}

func NewProtocolsRouter(kernel *core.Kernel, mux *http.ServeMux, middleware func(next http.Handler) http.Handler) *ProtocolsRouter {
	r := &ProtocolsRouter{
		kernel: kernel,
	}
	mux.Handle("GET /api/protocols", middleware(http.HandlerFunc(r.handleListProtocols)))
	return r
}

func (s *ProtocolsRouter) handleListProtocols(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListProtocols())
}

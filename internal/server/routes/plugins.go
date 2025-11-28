package routes

import (
	"encoding/json"
	"gohome/internal/core"
	"net/http"
)

type PluginsRouter struct {
	kernel *core.Kernel
}

func NewPluginsRouter(kernel *core.Kernel, mux *http.ServeMux, middleware func(next http.Handler) http.Handler) *PluginsRouter {
	r := &PluginsRouter{
		kernel: kernel,
	}
	mux.Handle("GET /api/plugins", middleware(http.HandlerFunc(r.handleListPlugins)))
	return r
}

func (s *PluginsRouter) handleListPlugins(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.kernel.ListPlugins())
}

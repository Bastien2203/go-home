package server

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Bastien2203/go-home/internal/core"
	"github.com/Bastien2203/go-home/internal/repository"
	"github.com/Bastien2203/go-home/internal/server/routes"
	"github.com/Bastien2203/go-home/internal/websockets"
	"github.com/Bastien2203/go-home/shared/config"
	"github.com/Bastien2203/go-home/shared/middlewares"
)

type Server struct {
	kernel         *core.Kernel
	addr           string
	wsHub          *websockets.Hub
	userRepository *repository.UserRepository
	sessionSecret  string
	appEnv         config.AppEnv
}

func NewServer(kernel *core.Kernel, port int, sessionSecret string, appEnv config.AppEnv, wsHub *websockets.Hub, userRepository *repository.UserRepository) *Server {
	return &Server{
		kernel:         kernel,
		addr:           fmt.Sprintf(":%d", port),
		wsHub:          wsHub,
		sessionSecret:  sessionSecret,
		appEnv:         appEnv,
		userRepository: userRepository,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticDir := "./dist"

	// --- Routes ---

	userRouter := routes.NewUsersRouter(mux, s.sessionSecret, s.appEnv, s.userRepository)
	routes.NewAdaptersRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewDevicesRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewPluginsRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewScannersRouter(s.kernel, mux, userRouter.AuthMiddleware)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(s.wsHub, w, r)
	})

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
		Handler: middlewares.CorsMiddleware(mux),
	}

	log.Printf("[Server] API listening on http://localhost%s", s.addr)
	return server.ListenAndServe()
}

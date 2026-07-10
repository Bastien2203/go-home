package server

import (
	"context"
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
	allowedOrigins []string
	httpServer     *http.Server
}

func NewServer(kernel *core.Kernel, port int, sessionSecret string, appEnv config.AppEnv, wsHub *websockets.Hub, userRepository *repository.UserRepository, allowedOrigins []string) *Server {
	return &Server{
		kernel:         kernel,
		addr:           fmt.Sprintf(":%d", port),
		wsHub:          wsHub,
		sessionSecret:  sessionSecret,
		appEnv:         appEnv,
		userRepository: userRepository,
		allowedOrigins: allowedOrigins,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticDir := "./dist"
	fs := http.FileServer(http.Dir(staticDir))

	// --- Routes ---

	userRouter := routes.NewUsersRouter(mux, s.sessionSecret, s.appEnv, s.userRepository)
	routes.NewAdaptersRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewDevicesRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewPluginsRouter(s.kernel, mux, userRouter.AuthMiddleware)
	routes.NewScannersRouter(s.kernel, mux, userRouter.AuthMiddleware)

	mux.Handle("/ws", userRouter.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(s.wsHub, w, r, s.allowedOrigins, s.appEnv)
	})))

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

		fs.ServeHTTP(w, r)
	})

	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: middlewares.CorsMiddleware(s.allowedOrigins, s.appEnv)(mux),
	}

	log.Printf("[Server] API listening on http://localhost%s", s.addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
